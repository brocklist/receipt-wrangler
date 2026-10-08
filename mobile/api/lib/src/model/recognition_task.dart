//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/recognition_task_stage.dart';
import 'package:openapi/src/model/recognition_task_status.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'recognition_task.g.dart';

/// RecognitionTask
///
/// Properties:
/// * [id] 
/// * [clientRequestId] 
/// * [fileName] 
/// * [fileSize] 
/// * [groupId] 
/// * [ownerUserId] 
/// * [version] 
/// * [status] 
/// * [stage] 
/// * [createdAt] 
/// * [updatedAt] 
/// * [queuedAt] 
/// * [startedAt] 
/// * [stageStartedAt] 
/// * [completedAt] 
/// * [uploadedBytes] 
/// * [uploadTotalBytes] 
/// * [attempt] 
/// * [maxAttempts] 
/// * [nextRetryAt] 
/// * [fallbackActive] 
/// * [receiptId] 
/// * [errorCode] 
/// * [errorMessage] 
/// * [canUpload] 
/// * [canRetry] 
@BuiltValue()
abstract class RecognitionTask implements Built<RecognitionTask, RecognitionTaskBuilder> {
  @BuiltValueField(wireName: r'id')
  int get id;

  @BuiltValueField(wireName: r'clientRequestId')
  String get clientRequestId;

  @BuiltValueField(wireName: r'fileName')
  String get fileName;

  @BuiltValueField(wireName: r'fileSize')
  int get fileSize;

  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  @BuiltValueField(wireName: r'ownerUserId')
  int get ownerUserId;

  @BuiltValueField(wireName: r'version')
  int get version;

  @BuiltValueField(wireName: r'status')
  RecognitionTaskStatus get status;
  // enum statusEnum {  AWAITING_UPLOAD,  UPLOADING,  UPLOAD_INTERRUPTED,  DISPATCH_PENDING,  QUEUED,  RUNNING,  RETRY_WAIT,  SUCCEEDED,  FAILED,  };

  @BuiltValueField(wireName: r'stage')
  RecognitionTaskStage get stage;
  // enum stageEnum {  UPLOAD,  PREPROCESSING,  OCR,  AI,  PARSING,  SAVING,  DONE,  };

  @BuiltValueField(wireName: r'createdAt')
  DateTime get createdAt;

  @BuiltValueField(wireName: r'updatedAt')
  DateTime get updatedAt;

  @BuiltValueField(wireName: r'queuedAt')
  DateTime? get queuedAt;

  @BuiltValueField(wireName: r'startedAt')
  DateTime? get startedAt;

  @BuiltValueField(wireName: r'stageStartedAt')
  DateTime? get stageStartedAt;

  @BuiltValueField(wireName: r'completedAt')
  DateTime? get completedAt;

  @BuiltValueField(wireName: r'uploadedBytes')
  int get uploadedBytes;

  @BuiltValueField(wireName: r'uploadTotalBytes')
  int? get uploadTotalBytes;

  @BuiltValueField(wireName: r'attempt')
  int get attempt;

  @BuiltValueField(wireName: r'maxAttempts')
  int get maxAttempts;

  @BuiltValueField(wireName: r'nextRetryAt')
  DateTime? get nextRetryAt;

  @BuiltValueField(wireName: r'fallbackActive')
  bool get fallbackActive;

  @BuiltValueField(wireName: r'receiptId')
  int? get receiptId;

  @BuiltValueField(wireName: r'errorCode')
  String get errorCode;

  @BuiltValueField(wireName: r'errorMessage')
  String get errorMessage;

  @BuiltValueField(wireName: r'canUpload')
  bool get canUpload;

  @BuiltValueField(wireName: r'canRetry')
  bool get canRetry;

  RecognitionTask._();

  factory RecognitionTask([void updates(RecognitionTaskBuilder b)]) = _$RecognitionTask;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(RecognitionTaskBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<RecognitionTask> get serializer => _$RecognitionTaskSerializer();
}

class _$RecognitionTaskSerializer implements PrimitiveSerializer<RecognitionTask> {
  @override
  final Iterable<Type> types = const [RecognitionTask, _$RecognitionTask];

  @override
  final String wireName = r'RecognitionTask';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    RecognitionTask object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'id';
    yield serializers.serialize(
      object.id,
      specifiedType: const FullType(int),
    );
    yield r'clientRequestId';
    yield serializers.serialize(
      object.clientRequestId,
      specifiedType: const FullType(String),
    );
    yield r'fileName';
    yield serializers.serialize(
      object.fileName,
      specifiedType: const FullType(String),
    );
    yield r'fileSize';
    yield serializers.serialize(
      object.fileSize,
      specifiedType: const FullType(int),
    );
    yield r'groupId';
    yield serializers.serialize(
      object.groupId,
      specifiedType: const FullType(int),
    );
    yield r'ownerUserId';
    yield serializers.serialize(
      object.ownerUserId,
      specifiedType: const FullType(int),
    );
    yield r'version';
    yield serializers.serialize(
      object.version,
      specifiedType: const FullType(int),
    );
    yield r'status';
    yield serializers.serialize(
      object.status,
      specifiedType: const FullType(RecognitionTaskStatus),
    );
    yield r'stage';
    yield serializers.serialize(
      object.stage,
      specifiedType: const FullType(RecognitionTaskStage),
    );
    yield r'createdAt';
    yield serializers.serialize(
      object.createdAt,
      specifiedType: const FullType(DateTime),
    );
    yield r'updatedAt';
    yield serializers.serialize(
      object.updatedAt,
      specifiedType: const FullType(DateTime),
    );
    if (object.queuedAt != null) {
      yield r'queuedAt';
      yield serializers.serialize(
        object.queuedAt,
        specifiedType: const FullType(DateTime),
      );
    }
    if (object.startedAt != null) {
      yield r'startedAt';
      yield serializers.serialize(
        object.startedAt,
        specifiedType: const FullType(DateTime),
      );
    }
    if (object.stageStartedAt != null) {
      yield r'stageStartedAt';
      yield serializers.serialize(
        object.stageStartedAt,
        specifiedType: const FullType(DateTime),
      );
    }
    if (object.completedAt != null) {
      yield r'completedAt';
      yield serializers.serialize(
        object.completedAt,
        specifiedType: const FullType(DateTime),
      );
    }
    yield r'uploadedBytes';
    yield serializers.serialize(
      object.uploadedBytes,
      specifiedType: const FullType(int),
    );
    if (object.uploadTotalBytes != null) {
      yield r'uploadTotalBytes';
      yield serializers.serialize(
        object.uploadTotalBytes,
        specifiedType: const FullType(int),
      );
    }
    yield r'attempt';
    yield serializers.serialize(
      object.attempt,
      specifiedType: const FullType(int),
    );
    yield r'maxAttempts';
    yield serializers.serialize(
      object.maxAttempts,
      specifiedType: const FullType(int),
    );
    if (object.nextRetryAt != null) {
      yield r'nextRetryAt';
      yield serializers.serialize(
        object.nextRetryAt,
        specifiedType: const FullType(DateTime),
      );
    }
    yield r'fallbackActive';
    yield serializers.serialize(
      object.fallbackActive,
      specifiedType: const FullType(bool),
    );
    if (object.receiptId != null) {
      yield r'receiptId';
      yield serializers.serialize(
        object.receiptId,
        specifiedType: const FullType(int),
      );
    }
    yield r'errorCode';
    yield serializers.serialize(
      object.errorCode,
      specifiedType: const FullType(String),
    );
    yield r'errorMessage';
    yield serializers.serialize(
      object.errorMessage,
      specifiedType: const FullType(String),
    );
    yield r'canUpload';
    yield serializers.serialize(
      object.canUpload,
      specifiedType: const FullType(bool),
    );
    yield r'canRetry';
    yield serializers.serialize(
      object.canRetry,
      specifiedType: const FullType(bool),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    RecognitionTask object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required RecognitionTaskBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'id':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.id = valueDes;
          break;
        case r'clientRequestId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.clientRequestId = valueDes;
          break;
        case r'fileName':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.fileName = valueDes;
          break;
        case r'fileSize':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.fileSize = valueDes;
          break;
        case r'groupId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.groupId = valueDes;
          break;
        case r'ownerUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.ownerUserId = valueDes;
          break;
        case r'version':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.version = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(RecognitionTaskStatus),
          ) as RecognitionTaskStatus;
          result.status = valueDes;
          break;
        case r'stage':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(RecognitionTaskStage),
          ) as RecognitionTaskStage;
          result.stage = valueDes;
          break;
        case r'createdAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.createdAt = valueDes;
          break;
        case r'updatedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.updatedAt = valueDes;
          break;
        case r'queuedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.queuedAt = valueDes;
          break;
        case r'startedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.startedAt = valueDes;
          break;
        case r'stageStartedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.stageStartedAt = valueDes;
          break;
        case r'completedAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.completedAt = valueDes;
          break;
        case r'uploadedBytes':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.uploadedBytes = valueDes;
          break;
        case r'uploadTotalBytes':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.uploadTotalBytes = valueDes;
          break;
        case r'attempt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.attempt = valueDes;
          break;
        case r'maxAttempts':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.maxAttempts = valueDes;
          break;
        case r'nextRetryAt':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(DateTime),
          ) as DateTime;
          result.nextRetryAt = valueDes;
          break;
        case r'fallbackActive':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.fallbackActive = valueDes;
          break;
        case r'receiptId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.receiptId = valueDes;
          break;
        case r'errorCode':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.errorCode = valueDes;
          break;
        case r'errorMessage':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(String),
          ) as String;
          result.errorMessage = valueDes;
          break;
        case r'canUpload':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.canUpload = valueDes;
          break;
        case r'canRetry':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(bool),
          ) as bool;
          result.canRetry = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  RecognitionTask deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = RecognitionTaskBuilder();
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

