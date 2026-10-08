//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'recognition_task_status.g.dart';

class RecognitionTaskStatus extends EnumClass {

  @BuiltValueEnumConst(wireName: r'AWAITING_UPLOAD')
  static const RecognitionTaskStatus AWAITING_UPLOAD = _$AWAITING_UPLOAD;
  @BuiltValueEnumConst(wireName: r'UPLOADING')
  static const RecognitionTaskStatus UPLOADING = _$UPLOADING;
  @BuiltValueEnumConst(wireName: r'UPLOAD_INTERRUPTED')
  static const RecognitionTaskStatus UPLOAD_INTERRUPTED = _$UPLOAD_INTERRUPTED;
  @BuiltValueEnumConst(wireName: r'DISPATCH_PENDING')
  static const RecognitionTaskStatus DISPATCH_PENDING = _$DISPATCH_PENDING;
  @BuiltValueEnumConst(wireName: r'QUEUED')
  static const RecognitionTaskStatus QUEUED = _$QUEUED;
  @BuiltValueEnumConst(wireName: r'RUNNING')
  static const RecognitionTaskStatus RUNNING = _$RUNNING;
  @BuiltValueEnumConst(wireName: r'RETRY_WAIT')
  static const RecognitionTaskStatus RETRY_WAIT = _$RETRY_WAIT;
  @BuiltValueEnumConst(wireName: r'SUCCEEDED')
  static const RecognitionTaskStatus SUCCEEDED = _$SUCCEEDED;
  @BuiltValueEnumConst(wireName: r'FAILED')
  static const RecognitionTaskStatus FAILED = _$FAILED;

  static Serializer<RecognitionTaskStatus> get serializer => _$recognitionTaskStatusSerializer;

  const RecognitionTaskStatus._(String name): super(name);

  static BuiltSet<RecognitionTaskStatus> get values => _$values;
  static RecognitionTaskStatus valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class RecognitionTaskStatusMixin = Object with _$RecognitionTaskStatusMixin;

