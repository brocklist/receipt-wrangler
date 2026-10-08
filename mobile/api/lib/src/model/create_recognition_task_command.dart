//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:openapi/src/model/receipt_status.dart';
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'create_recognition_task_command.g.dart';

/// CreateRecognitionTaskCommand
///
/// Properties:
/// * [clientRequestId] 
/// * [fileName] 
/// * [fileSize] 
/// * [groupId] 
/// * [paidByUserId] 
/// * [status] 
/// * [categoryIds] 
/// * [tagIds] 
@BuiltValue()
abstract class CreateRecognitionTaskCommand implements Built<CreateRecognitionTaskCommand, CreateRecognitionTaskCommandBuilder> {
  @BuiltValueField(wireName: r'clientRequestId')
  String get clientRequestId;

  @BuiltValueField(wireName: r'fileName')
  String get fileName;

  @BuiltValueField(wireName: r'fileSize')
  int get fileSize;

  @BuiltValueField(wireName: r'groupId')
  int get groupId;

  @BuiltValueField(wireName: r'paidByUserId')
  int? get paidByUserId;

  @BuiltValueField(wireName: r'status')
  ReceiptStatus? get status;
  // enum statusEnum {  OPEN,  NEEDS_ATTENTION,  RESOLVED,  DRAFT,  ,  };

  @BuiltValueField(wireName: r'categoryIds')
  BuiltList<int>? get categoryIds;

  @BuiltValueField(wireName: r'tagIds')
  BuiltList<int>? get tagIds;

  CreateRecognitionTaskCommand._();

  factory CreateRecognitionTaskCommand([void updates(CreateRecognitionTaskCommandBuilder b)]) = _$CreateRecognitionTaskCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(CreateRecognitionTaskCommandBuilder b) => b
      ..paidByUserId = 0;

  @BuiltValueSerializer(custom: true)
  static Serializer<CreateRecognitionTaskCommand> get serializer => _$CreateRecognitionTaskCommandSerializer();
}

class _$CreateRecognitionTaskCommandSerializer implements PrimitiveSerializer<CreateRecognitionTaskCommand> {
  @override
  final Iterable<Type> types = const [CreateRecognitionTaskCommand, _$CreateRecognitionTaskCommand];

  @override
  final String wireName = r'CreateRecognitionTaskCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    CreateRecognitionTaskCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
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
    if (object.paidByUserId != null) {
      yield r'paidByUserId';
      yield serializers.serialize(
        object.paidByUserId,
        specifiedType: const FullType(int),
      );
    }
    if (object.status != null) {
      yield r'status';
      yield serializers.serialize(
        object.status,
        specifiedType: const FullType(ReceiptStatus),
      );
    }
    if (object.categoryIds != null) {
      yield r'categoryIds';
      yield serializers.serialize(
        object.categoryIds,
        specifiedType: const FullType(BuiltList, [FullType(int)]),
      );
    }
    if (object.tagIds != null) {
      yield r'tagIds';
      yield serializers.serialize(
        object.tagIds,
        specifiedType: const FullType(BuiltList, [FullType(int)]),
      );
    }
  }

  @override
  Object serialize(
    Serializers serializers,
    CreateRecognitionTaskCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required CreateRecognitionTaskCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
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
        case r'paidByUserId':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.paidByUserId = valueDes;
          break;
        case r'status':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(ReceiptStatus),
          ) as ReceiptStatus;
          result.status = valueDes;
          break;
        case r'categoryIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.categoryIds.replace(valueDes);
          break;
        case r'tagIds':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(BuiltList, [FullType(int)]),
          ) as BuiltList<int>;
          result.tagIds.replace(valueDes);
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  CreateRecognitionTaskCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = CreateRecognitionTaskCommandBuilder();
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

