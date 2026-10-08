//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/recognition_task.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'get_recognition_tasks_response.g.dart';

/// GetRecognitionTasksResponse
///
/// Properties:
/// * [data] 
/// * [totalCount] 
/// * [activeCount] 
/// * [awaitingUploadCount] 
/// * [runningCount] 
/// * [failedCount] 
@BuiltValue()
abstract class GetRecognitionTasksResponse implements Built<GetRecognitionTasksResponse, GetRecognitionTasksResponseBuilder> {
  @BuiltValueField(wireName: r'data')
  BuiltList<RecognitionTask> get data;

  @BuiltValueField(wireName: r'totalCount')
  int get totalCount;

  @BuiltValueField(wireName: r'activeCount')
  int get activeCount;

  @BuiltValueField(wireName: r'awaitingUploadCount')
  int get awaitingUploadCount;

  @BuiltValueField(wireName: r'runningCount')
  int get runningCount;

  @BuiltValueField(wireName: r'failedCount')
  int get failedCount;

  GetRecognitionTasksResponse._();

  factory GetRecognitionTasksResponse([void updates(GetRecognitionTasksResponseBuilder b)]) = _$GetRecognitionTasksResponse;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(GetRecognitionTasksResponseBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<GetRecognitionTasksResponse> get serializer => _$GetRecognitionTasksResponseSerializer();
}

class _$GetRecognitionTasksResponseSerializer implements PrimitiveSerializer<GetRecognitionTasksResponse> {
  @override
  final Iterable<Type> types = const [GetRecognitionTasksResponse, _$GetRecognitionTasksResponse];

  @override
  final String wireName = r'GetRecognitionTasksResponse';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    GetRecognitionTasksResponse object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'data';
    yield serializers.serialize(
      object.data,
      specifiedType: const FullType(BuiltList, [FullType(RecognitionTask)]),
    );
    yield r'totalCount';
    yield serializers.serialize(
      object.totalCount,
      specifiedType: const FullType(int),
    );
    yield r'activeCount';
    yield serializers.serialize(
      object.activeCount,
      specifiedType: const FullType(int),
    );
    yield r'awaitingUploadCount';
    yield serializers.serialize(
      object.awaitingUploadCount,
      specifiedType: const FullType(int),
    );
    yield r'runningCount';
    yield serializers.serialize(
      object.runningCount,
      specifiedType: const FullType(int),
    );
    yield r'failedCount';
    yield serializers.serialize(
      object.failedCount,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    GetRecognitionTasksResponse object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required GetRecognitionTasksResponseBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'data':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(RecognitionTask)]),
          ) as BuiltList<RecognitionTask>;
          result.data.replace(valueDes);
          break;
        case r'totalCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.totalCount = valueDes;
          break;
        case r'activeCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.activeCount = valueDes;
          break;
        case r'awaitingUploadCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.awaitingUploadCount = valueDes;
          break;
        case r'runningCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.runningCount = valueDes;
          break;
        case r'failedCount':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.failedCount = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  GetRecognitionTasksResponse deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = GetRecognitionTasksResponseBuilder();
    final serializedList = (serialized as Iterable<Object?>).toList();
    final unhandled = <Object?>[];
    _deserializeProperties(
      serializers,
      serialized,
      specifiedType: specifiedType,
      serializedList: serializedList,
      unhandled: unhandled,
      result: result,
    );
    return result.build();
  }
}

