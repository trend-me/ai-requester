Feature: Build AI prompts and forward the event to ai-requester queue

  Scenario: Successfully process a message from the queue
    Given a message with the following data is sent to 'ai-requester' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "prompt_road_map_step":2,
    "model":"gemini",
    "prompt":"this is a test. [1 2 3 4] 1",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given the prompt road map API returns the following prompt road map for step '2' and prompt_road_map_config_name 'TEST':
    """
    {
    "prompt_road_map_config_name":"TEST",
    "response_validation_name":"TEST_RESPONSE",
    "metadata_validation_name":"TEST_METADATA",
    "question_template":"this is a <any.thing>. <any.array> <any.array[0]>",
    "step":2,
    "created_at":"2024-08-01T20:53:49.132Z",
    "updated_at":"2024-08-01T20:53:49.132Z"
    }
    """
    Given the validation API returns the following validation result for payload_validation 'TEST_RESPONSE':
    """
    {
      "failures": "",
      "errors":[]
    }
    """
    Given the ai model 'gemini' returns the following response:
    """
    {
      "response":"valid"
    }
    """
    When the message is consumed by the ai-requester consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name 'TEST' and step '2'
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name 'TEST' and step '2'
    And the following prompt should be sent to the following ai model:
    |         prompt             |   model  |
    | this is a test. [1 2 3 4] 1 |   gemini |
    And the response should be sent to the validation API with the payload_validation 'TEST_RESPONSE'
    And a message with the following data should be sent to 'ai-callback' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "model":"gemini",
    "prompt_road_map_step":2,
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]}, "response":"valid" } 
    }
    """
    And the application should not retry


   Scenario: Successfully handles a response validation failure
    Given a message with the following data is sent to 'ai-requester' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "prompt":"this is a test. [1 2 3 4] 1",
    "prompt_road_map_step":2,
    "model":"gemini",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given the prompt road map API returns the following prompt road map for step '2' and prompt_road_map_config_name 'TEST':
    """
    {
    "prompt_road_map_config_name":"TEST",
    "response_validation_name":"TEST_RESPONSE",
    "metadata_validation_name":"TEST_METADATA",
    "question_template":"this is a <any.thing>. <any.array> <any.array[0]>",
    "step":2,
    "created_at":"2024-08-01T20:53:49.132Z",
    "updated_at":"2024-08-01T20:53:49.132Z"
    }
    """
    Given the ai model 'gemini' returns the following response:
    """
    {
      "response":"invalid"
    }
    """
    Given the validation API returns the following validation result for payload_validation 'TEST_RESPONSE':
    """
    {
      "failures": "there is something wrong",
      "errors":[]
    }
    """
    When the message is consumed by the ai-requester consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name 'TEST' and step '2'
    And the following prompt should be sent to the following ai model:
    |         prompt              |   model  |
    | this is a test. [1 2 3 4] 1 |   gemini |
    And the response should be sent to the validation API with the payload_validation 'TEST_RESPONSE'
    And no message should be sent to the 'ai-callback' queue
    And the application should retry

  Scenario: Successfully process an error and scheduling a retry
    Given a message with the following data is sent to 'ai-requester' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "prompt_road_map_step":2,
    "prompt":"this is a test. [1 2 3 4] 1",
    "model":"gemini",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given the prompt road map API returns an statusCode 500
    When the message is consumed by the ai-requester consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name 'TEST' and step '2'
    And the response should not be sent to the validation API
    And no message should be sent to the 'ai-callback' queue
    And the application should retry

  Scenario: Successfully retries on a gemini error
    Given a message with the following data is sent to 'ai-requester' queue:
    """
    {
    "prompt_road_map_config_execution_id":"c713deb9-efa2-4d5f-9675-abe0b7e0c0d4",
    "prompt_road_map_config_name":"TEST",
    "output_queue":"output-queue",
    "prompt_road_map_step":2,
    "prompt":"this is a test. [1 2 3 4] 1",
    "model":"gemini",
    "metadata":{"any": { "thing":"test", "array":[1,2,3,4]} }
    }
    """
    Given the ai model 'gemini' fails with an error 'some error'
    Given the prompt road map API returns the following prompt road map for step '2' and prompt_road_map_config_name 'TEST':
    """
    {
    "prompt_road_map_config_name":"TEST",
    "response_validation_name":"TEST_RESPONSE",
    "metadata_validation_name":"TEST_METADATA",
    "question_template":"this is a <any.thing>. <any.array> <any.array[0]>",
    "step":2,
    "created_at":"2024-08-01T20:53:49.132Z",
    "updated_at":"2024-08-01T20:53:49.132Z"
    }
    """
    When the message is consumed by the ai-requester consumer
    Then the prompt_road_map is fetched from the prompt-road-map-api using the prompt_road_map_config_name 'TEST' and step '2'
    And the following prompt should be sent to the following ai model:
    |         prompt              |   model  |
    | this is a test. [1 2 3 4] 1 |   gemini |
    And the response should not be sent to the validation API
    And no message should be sent to the 'ai-callback' queue
    And the application should retry
