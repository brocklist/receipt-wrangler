# openapi.model.RecognitionTask

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**clientRequestId** | **String** |  | 
**fileName** | **String** |  | 
**fileSize** | **int** |  | 
**groupId** | **int** |  | 
**ownerUserId** | **int** |  | 
**version** | **int** |  | 
**status** | [**RecognitionTaskStatus**](RecognitionTaskStatus.md) |  | 
**stage** | [**RecognitionTaskStage**](RecognitionTaskStage.md) |  | 
**createdAt** | [**DateTime**](DateTime.md) |  | 
**updatedAt** | [**DateTime**](DateTime.md) |  | 
**queuedAt** | [**DateTime**](DateTime.md) |  | [optional] 
**startedAt** | [**DateTime**](DateTime.md) |  | [optional] 
**stageStartedAt** | [**DateTime**](DateTime.md) |  | [optional] 
**completedAt** | [**DateTime**](DateTime.md) |  | [optional] 
**uploadedBytes** | **int** |  | 
**uploadTotalBytes** | **int** |  | [optional] 
**attempt** | **int** |  | 
**maxAttempts** | **int** |  | 
**nextRetryAt** | [**DateTime**](DateTime.md) |  | [optional] 
**fallbackActive** | **bool** |  | 
**receiptId** | **int** |  | [optional] 
**errorCode** | **String** |  | 
**errorMessage** | **String** |  | 
**canUpload** | **bool** |  | 
**canRetry** | **bool** |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


