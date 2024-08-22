package steps

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/trend-me/ai-requester/internal/domain/interfaces"
	"github.com/trend-me/ai-requester/internal/integration/ais"

	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/trend-me/ai-requester/internal/config/properties"
	"github.com/trend-me/ai-requester/test/bdd/containers"
	rabbitmq_container "github.com/trend-me/ai-requester/test/bdd/containers/rabbitmq"
	"github.com/trend-me/ai-requester/test/bdd/injector"
	"github.com/trend-me/ai-requester/test/bdd/mocks"
	"github.com/trend-me/ai-requester/test/bdd/utils"
	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/params"
	"github.com/vitorsalgado/mocha/v3/reply"
)

var (
	t                                              *testing.T
	consumedMessage                                string
	consumer                                       interfaces.QueueAiRequesterConsumer
	m                                              *mocha.Mocha
	scopePromptRoadMapConfigsApiGetPromptRoadMap   *mocha.Scoped
	scopePayloadValidationApiExecute               *mocha.Scoped
	requestPayloadValidationApiExecute             *http.Request
	requestPromptRoadMapConfigsApiGetPromptRoadMap *http.Request
	geminiMock                                     *mocks.GeminiModelMock
)

func setup(t *testing.T) {
	m = mocha.New(t)
	m.Start()
	_ = os.Setenv("URL_API_PROMPT_ROAD_MAP_CONFIG", m.URL()+"/prompt_road_map_configs")
	_ = os.Setenv("URL_API_PROMPT_ROAD_MAP_CONFIG_EXECUTION", m.URL()+"/prompt_road_map_config_executions")
	_ = os.Setenv("URL_API_VALIDATION", m.URL()+"/payload_validations")
	err := godotenv.Load("../.bdd.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	containers.Up()
	for range 10 {
		time.Sleep(10 * time.Second)
		fmt.Println("Waiting for rabbitmq to start")
		err = rabbitmq_container.Connect()
		if err == nil {
			break
		}
		fmt.Println(err.Error())
	}
	if err != nil {
		t.Fatal(err.Error())
	}

	geminiMock = &mocks.GeminiModelMock{}
	consumer, err = injector.InitializeQueueAiRequesterConsumerMock(
		ais.NewGemini(
			func() []string {
				return []string{"test"}
			},
			func(ctx context.Context, key string) (ais.GeminiModel, error) {
				return geminiMock, nil
			},
		))
	if err != nil {
		t.Fatal(err.Error())
	}

}

func down(t *testing.T) {
	err := rabbitmq_container.Disconnect()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = containers.Down()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = m.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFeatures(t_ *testing.T) {
	t = t_
	setup(t)
	t.Cleanup(func() {
		defer down(t)
	})

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func aMessageWithTheFollowingDataIsSentToAipromptbuilderQueue(queue string, arg1 *godog.DocString) error {
	if queue == properties.QueueAiRequester {
		consumedMessage = arg1.Content
	}
	return rabbitmq_container.PostMessageToQueue(queue, []byte(arg1.Content))
}

func aMessageWithTheFollowingDataShouldBeSentToAipromptbuilderQueue(queue string, arg1 *godog.DocString) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(queue)
	if err != nil {
		return err
	}
	res, err := utils.JsonEqual(arg1.Content, string(content))
	if !res {
		if err != nil {
			return fmt.Errorf("error comparing json: %v", err)
		}
		return fmt.Errorf("message sent to queue '%s' is not equal to the expected message: %s. Got: %s",
			queue, arg1.Content, string(content))
	}

	return nil
}

func noMessageShouldBeSentToTheAirequesterQueue(queue string) error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(queue)
	if err != nil {
		return err
	}

	if !assert.Nil(t, content) {
		return fmt.Errorf("a message was sent to queue '%s'. Got: %s",
			queue, string(content))
	}
	return nil
}

func noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi() error {
	if scopePromptRoadMapConfigsApiGetPromptRoadMap != nil && scopePromptRoadMapConfigsApiGetPromptRoadMap.Called() {
		return fmt.Errorf("prompt road map was fetched")
	}

	return nil
}

func theApplicationShouldNotRetry() error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(properties.QueueAiRequester)
	if err != nil {
		return err
	}
	assert.Nil(t, content)
	return nil
}

func theApplicationShouldRetry() error {
	content, _, err := rabbitmq_container.ConsumeMessageFromQueue(properties.QueueAiRequester)
	if err != nil {
		return err
	}

	res, err := utils.JsonEqual(string(content), string(consumedMessage))
	if !res {
		if err != nil {
			return fmt.Errorf("error comparing json: %v", err)
		}

		return fmt.Errorf("message sent to queue '%s' is not equal to the expected message: %s. Got: %s",
			properties.QueueAiRequester, consumedMessage, string(content))

	}
	return nil
}

func theMessageIsConsumedByTheAipromptbuilderConsumer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh, err := consumer.Consume(ctx)
	if err != nil {
		return err
	}

	timeout := time.After(120 * time.Second)
	for {
		select {
		case <-errCh:
			cancel()
			return nil
		case <-timeout:
			err = fmt.Errorf("timeout")
			return err
		}
	}
}

func theresponseShouldBeSentToTheValidationAPIWithTheresponse_validation_nameTEST_response(name string) error {
	if !scopePayloadValidationApiExecute.Called() || requestPayloadValidationApiExecute == nil {
		return fmt.Errorf("response was not sent to the validation API")
	}

	split := strings.Split(requestPayloadValidationApiExecute.URL.Path, "/")
	if split[len(split)-1] != name {
		return fmt.Errorf("response was not sent to the validation API with the correct response_validation_name. I was %s",
			requestPayloadValidationApiExecute.URL.Path)
	}

	return nil
}

func theresponseShouldNotBeSentToTheValidationAPI() error {
	if scopePayloadValidationApiExecute != nil && scopePayloadValidationApiExecute.Called() {
		return fmt.Errorf("response was sent to the validation API")
	}

	return nil
}

func thePromptRoadMapAPIReturnsAnStatusCode500() error {
	scopePromptRoadMapConfigsApiGetPromptRoadMap = m.AddMocks(mocha.Get(expect.Func(func(v any, a expect.Args) (bool, error) {
		return strings.Contains(a.RequestInfo.Request.URL.Path, "/prompt_road_map_configs"), nil
	})).ReplyFunction(func(r *http.Request, _ reply.M, _ params.P) (*reply.Response, error) {
		requestPromptRoadMapConfigsApiGetPromptRoadMap = r
		return &reply.Response{
			Status: http.StatusInternalServerError,
			Body:   io.NopCloser(strings.NewReader(`{"error": "Internal Server Error"}`)),
		}, nil
	}))

	return nil
}

func thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap(step int, name string, arg1 *godog.DocString) error {
	scopePromptRoadMapConfigsApiGetPromptRoadMap = m.AddMocks(mocha.
		Get(expect.URLPath(fmt.Sprintf("/prompt_road_map_configs/%s/prompt_road_maps/%d", name, step))).ReplyFunction(func(request *http.Request, r reply.M, p params.P) (*reply.Response, error) {
		requestPromptRoadMapConfigsApiGetPromptRoadMap = request
		return &reply.Response{
			Status: http.StatusOK,
			Body:   io.NopCloser(strings.NewReader(arg1.Content)),
		}, nil
	}))

	return nil
}

func thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name(name string, step int) error {
	if scopePromptRoadMapConfigsApiGetPromptRoadMap == nil || !scopePromptRoadMapConfigsApiGetPromptRoadMap.Called() || requestPromptRoadMapConfigsApiGetPromptRoadMap == nil {
		return fmt.Errorf("prompt road map was not fetched")
	}

	if !strings.Contains(requestPromptRoadMapConfigsApiGetPromptRoadMap.URL.Path, fmt.Sprintf("/prompt_road_map_configs/%s/prompt_road_maps/%d", name, step)) {
		return fmt.Errorf("prompt_road_map_config fetched with '%s'. Required prompt_road_map_config_name: '%s' and step: '%d'",
			requestPromptRoadMapConfigsApiGetPromptRoadMap.URL.Path, name, step)
	}

	return nil
}

func theValidationAPIReturnsTheFollowingValidationResult(name string, arg1 *godog.DocString) error {
	scopePayloadValidationApiExecute = m.AddMocks(mocha.
		Post(expect.URLPath(fmt.Sprintf("/payload_validations/%s", name))).
		ReplyFunction(func(request *http.Request, r reply.M, p params.P) (*reply.Response, error) {
			requestPayloadValidationApiExecute = request
			return &reply.Response{
				Status: http.StatusOK,
				Body:   io.NopCloser(strings.NewReader(arg1.Content)),
			}, nil
		}))

	return nil
}

func theFollowingPromptShouldBeSentToTheFollowingAiModel(arg1 *godog.Table) error {
	if len(arg1.Rows[0].Cells) != 2 || arg1.Rows[0].Cells[0].Value != "prompt" || arg1.Rows[0].Cells[1].Value != "model" {
		return fmt.Errorf("expected table with headers 'prompt' and 'model'")
	}
	prompt := arg1.Rows[1].Cells[0].Value
	model := arg1.Rows[1].Cells[1].Value

	if model == properties.AiModelNameGemini {
		if geminiMock.GetPrompt() == prompt {
			return nil
		}
		return fmt.Errorf("expected prompt '%s' to be sent to model '%s'. Got: '%s'", prompt, model, geminiMock.GetPrompt())
	}

	return fmt.Errorf("expected prompt '%s' to be sent to model '%s'. But model was not configured in test setup", prompt, model)
}

func theAiModelReturnsTheFollowingResponse(model string, arg1 *godog.DocString) error {
	if model == properties.AiModelNameGemini {
		geminiMock.SetResponse(arg1.Content)
		return nil
	}
	return fmt.Errorf("model '%s' was not configured in test setup", model)
}

func theAiModelFailsWithAnError(model, error string) error {
	if model == properties.AiModelNameGemini {
		geminiMock.SetError(fmt.Errorf(error))
		return nil
	}
	return fmt.Errorf("model '%s' was not configured in test setup", model)
}

func maxReceiveCountIs(count string) error {
	os.Setenv("MAX_RECEIVE_COUNT", count)
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		geminiMock.Clean()

		if scopePromptRoadMapConfigsApiGetPromptRoadMap != nil {
			scopePromptRoadMapConfigsApiGetPromptRoadMap.Clean()
			scopePromptRoadMapConfigsApiGetPromptRoadMap = nil
		}

		if scopePayloadValidationApiExecute != nil {
			scopePayloadValidationApiExecute.Clean()
			scopePayloadValidationApiExecute = nil
		}

		scopePayloadValidationApiExecute = nil
		requestPayloadValidationApiExecute = nil
		requestPromptRoadMapConfigsApiGetPromptRoadMap = nil
		consumedMessage = ""

		_ = rabbitmq_container.PurgeMessages()
		return ctx, nil
	})

	ctx.Step(`^a message with the following data is sent to \'(.*)\' queue:$`, aMessageWithTheFollowingDataIsSentToAipromptbuilderQueue)
	ctx.Step(`^a message with the following data should be sent to \'(.*)\' queue:$`, aMessageWithTheFollowingDataShouldBeSentToAipromptbuilderQueue)
	ctx.Step(`^no message should be sent to the \'(.*)\' queue$`, noMessageShouldBeSentToTheAirequesterQueue)
	ctx.Step(`^no prompt_road_map should be fetched from the prompt-road-map-api$`, noPrompt_road_mapShouldBeFetchedFromThePromptroadmapapi)
	ctx.Step(`^the application should not retry$`, theApplicationShouldNotRetry)
	ctx.Step(`^the application should retry$`, theApplicationShouldRetry)
	ctx.Step(`^max receive count is \'(.*)\'$`, maxReceiveCountIs)
	ctx.Step(`^the message is consumed by the ai-requester consumer$`, theMessageIsConsumedByTheAipromptbuilderConsumer)
	ctx.Step(`^the response should be sent to the validation API with the payload_validation \'(.*)\'$`, theresponseShouldBeSentToTheValidationAPIWithTheresponse_validation_nameTEST_response)
	ctx.Step(`^the response should not be sent to the validation API$`, theresponseShouldNotBeSentToTheValidationAPI)
	ctx.Step(`^the prompt road map API returns an statusCode 500$`, thePromptRoadMapAPIReturnsAnStatusCode500)
	ctx.Step(`^the prompt road map API returns the following prompt road map for step \'(\d+)\' and prompt_road_map_config_name \'(.*)\':$`, thePromptRoadMapAPIReturnsTheFollowingPromptRoadMap)
	ctx.Step(`^the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name \'(.*)\' and step \'(\d+)\'$`, thePrompt_road_mapIsFetchedFromThePromptroadmapapiUsingThePrompt_road_map_config_name)
	ctx.Step(`^the validation API returns the following validation result for payload_validation \'(.*)\':$`, theValidationAPIReturnsTheFollowingValidationResult)
	ctx.Step(`^the following prompt should be sent to the following ai model:$`, theFollowingPromptShouldBeSentToTheFollowingAiModel)
	ctx.Given(`the ai model \'(.*)\' returns the following response:`, theAiModelReturnsTheFollowingResponse)
	ctx.Given(`the ai model \'(.*)\' fails with an error \'(.*)\'`, theAiModelFailsWithAnError)
}
