# openapi.api.RecognitionTaskApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createRecognitionTask**](RecognitionTaskApi.md#createrecognitiontask) | **POST** /recognitionTask | Register a Quick Scan file before uploading
[**getRecognitionTask**](RecognitionTaskApi.md#getrecognitiontask) | **GET** /recognitionTask/{id} | Refresh one authorized task
[**getRecognitionTasks**](RecognitionTaskApi.md#getrecognitiontasks) | **GET** /recognitionTask | Get authorized Quick Scan tasks and scoped counts
[**retryRecognitionTask**](RecognitionTaskApi.md#retryrecognitiontask) | **POST** /recognitionTask/{id}/retry | Retry a failed task using its retained source
[**uploadRecognitionTaskFile**](RecognitionTaskApi.md#uploadrecognitiontaskfile) | **PUT** /recognitionTask/{id}/file | Upload a complete Quick Scan source file


# **createRecognitionTask**
> RecognitionTask createRecognitionTask(createRecognitionTaskCommand)

Register a Quick Scan file before uploading

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRecognitionTaskApi();
final CreateRecognitionTaskCommand createRecognitionTaskCommand = ; // CreateRecognitionTaskCommand | 

try {
    final response = api.createRecognitionTask(createRecognitionTaskCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RecognitionTaskApi->createRecognitionTask: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **createRecognitionTaskCommand** | [**CreateRecognitionTaskCommand**](CreateRecognitionTaskCommand.md)|  | 

### Return type

[**RecognitionTask**](RecognitionTask.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRecognitionTask**
> RecognitionTask getRecognitionTask(id)

Refresh one authorized task

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRecognitionTaskApi();
final int id = 56; // int | 

try {
    final response = api.getRecognitionTask(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RecognitionTaskApi->getRecognitionTask: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**|  | 

### Return type

[**RecognitionTask**](RecognitionTask.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRecognitionTasks**
> GetRecognitionTasksResponse getRecognitionTasks(scope, bucket, page, pageSize, clientRequestId, ids, groupId, status)

Get authorized Quick Scan tasks and scoped counts

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRecognitionTaskApi();
final String scope = scope_example; // String | 
final String bucket = bucket_example; // String | 
final int page = 56; // int | 
final int pageSize = 56; // int | 
final String clientRequestId = clientRequestId_example; // String | 
final String ids = ids_example; // String | Comma-separated task IDs (at most 100)
final int groupId = 56; // int | 
final RecognitionTaskStatus status = ; // RecognitionTaskStatus | 

try {
    final response = api.getRecognitionTasks(scope, bucket, page, pageSize, clientRequestId, ids, groupId, status);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RecognitionTaskApi->getRecognitionTasks: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | **String**|  | [optional] [default to 'own']
 **bucket** | **String**|  | [optional] [default to 'all']
 **page** | **int**|  | [optional] [default to 1]
 **pageSize** | **int**|  | [optional] [default to 25]
 **clientRequestId** | **String**|  | [optional] 
 **ids** | **String**| Comma-separated task IDs (at most 100) | [optional] 
 **groupId** | **int**|  | [optional] 
 **status** | [**RecognitionTaskStatus**](.md)|  | [optional] 

### Return type

[**GetRecognitionTasksResponse**](GetRecognitionTasksResponse.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **retryRecognitionTask**
> RecognitionTask retryRecognitionTask(id, retryRecognitionTaskCommand)

Retry a failed task using its retained source

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRecognitionTaskApi();
final int id = 56; // int | 
final RetryRecognitionTaskCommand retryRecognitionTaskCommand = ; // RetryRecognitionTaskCommand | 

try {
    final response = api.retryRecognitionTask(id, retryRecognitionTaskCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RecognitionTaskApi->retryRecognitionTask: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**|  | 
 **retryRecognitionTaskCommand** | [**RetryRecognitionTaskCommand**](RetryRecognitionTaskCommand.md)|  | 

### Return type

[**RecognitionTask**](RecognitionTask.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uploadRecognitionTaskFile**
> RecognitionTask uploadRecognitionTaskFile(id, file)

Upload a complete Quick Scan source file

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRecognitionTaskApi();
final int id = 56; // int | 
final MultipartFile file = BINARY_DATA_HERE; // MultipartFile | 

try {
    final response = api.uploadRecognitionTaskFile(id, file);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RecognitionTaskApi->uploadRecognitionTaskFile: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**|  | 
 **file** | **MultipartFile**|  | 

### Return type

[**RecognitionTask**](RecognitionTask.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

