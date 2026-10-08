# openapi.model.CreateRecognitionTaskCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**clientRequestId** | **String** |  | 
**fileName** | **String** |  | 
**fileSize** | **int** |  | 
**groupId** | **int** |  | 
**paidByUserId** | **int** |  | [optional] [default to 0]
**status** | [**ReceiptStatus**](ReceiptStatus.md) |  | [optional] 
**categoryIds** | **BuiltList&lt;int&gt;** |  | [optional] 
**tagIds** | **BuiltList&lt;int&gt;** |  | [optional] 
**comment** | **String** | Optional quick scan receipt comment; group and role rules may require it. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


