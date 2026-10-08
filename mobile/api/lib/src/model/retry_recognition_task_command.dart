//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'retry_recognition_task_command.g.dart';

/// RetryRecognitionTaskCommand
///
/// Properties:
/// * [version] 
@BuiltValue()
abstract class RetryRecognitionTaskCommand implements Built<RetryRecognitionTaskCommand, RetryRecognitionTaskCommandBuilder> {
  @BuiltValueField(wireName: r'version')
  int get version;

  RetryRecognitionTaskCommand._();

  factory RetryRecognitionTaskCommand([void updates(RetryRecognitionTaskCommandBuilder b)]) = _$RetryRecognitionTaskCommand;

  @BuiltValueHook(initializeBuilder: true)
  static void _defaults(RetryRecognitionTaskCommandBuilder b) => b;

  @BuiltValueSerializer(custom: true)
  static Serializer<RetryRecognitionTaskCommand> get serializer => _$RetryRecognitionTaskCommandSerializer();
}

class _$RetryRecognitionTaskCommandSerializer implements PrimitiveSerializer<RetryRecognitionTaskCommand> {
  @override
  final Iterable<Type> types = const [RetryRecognitionTaskCommand, _$RetryRecognitionTaskCommand];

  @override
  final String wireName = r'RetryRecognitionTaskCommand';

  Iterable<Object?> _serializeProperties(
    Serializers serializers,
    RetryRecognitionTaskCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) sync* {
    yield r'version';
    yield serializers.serialize(
      object.version,
      specifiedType: const FullType(int),
    );
  }

  @override
  Object serialize(
    Serializers serializers,
    RetryRecognitionTaskCommand object, {
    FullType specifiedType = FullType.unspecified,
  }) {
    return _serializeProperties(serializers, object, specifiedType: specifiedType).toList();
  }

  void _deserializeProperties(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
    required List<Object?> serializedList,
    required RetryRecognitionTaskCommandBuilder result,
    required List<Object?> unhandled,
  }) {
    for (var i = 0; i < serializedList.length; i += 2) {
      final key = serializedList[i] as String;
      final value = serializedList[i + 1];
      switch (key) {
        case r'version':
          final valueDes = serializers.deserialize(
            value,
            specifiedType: const FullType(int),
          ) as int;
          result.version = valueDes;
          break;
        default:
          unhandled.add(key);
          unhandled.add(value);
          break;
      }
    }
  }

  @override
  RetryRecognitionTaskCommand deserialize(
    Serializers serializers,
    Object serialized, {
    FullType specifiedType = FullType.unspecified,
  }) {
    final result = RetryRecognitionTaskCommandBuilder();
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

