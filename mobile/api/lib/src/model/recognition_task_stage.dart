//
// AUTO-GENERATED FILE, DO NOT MODIFY!
//

// ignore_for_file: unused_element
import 'package:built_collection/built_collection.dart';
import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'recognition_task_stage.g.dart';

class RecognitionTaskStage extends EnumClass {

  @BuiltValueEnumConst(wireName: r'UPLOAD')
  static const RecognitionTaskStage UPLOAD = _$UPLOAD;
  @BuiltValueEnumConst(wireName: r'PREPROCESSING')
  static const RecognitionTaskStage PREPROCESSING = _$PREPROCESSING;
  @BuiltValueEnumConst(wireName: r'OCR')
  static const RecognitionTaskStage OCR = _$OCR;
  @BuiltValueEnumConst(wireName: r'AI')
  static const RecognitionTaskStage AI = _$AI;
  @BuiltValueEnumConst(wireName: r'PARSING')
  static const RecognitionTaskStage PARSING = _$PARSING;
  @BuiltValueEnumConst(wireName: r'SAVING')
  static const RecognitionTaskStage SAVING = _$SAVING;
  @BuiltValueEnumConst(wireName: r'DONE')
  static const RecognitionTaskStage DONE = _$DONE;

  static Serializer<RecognitionTaskStage> get serializer => _$recognitionTaskStageSerializer;

  const RecognitionTaskStage._(String name): super(name);

  static BuiltSet<RecognitionTaskStage> get values => _$values;
  static RecognitionTaskStage valueOf(String name) => _$valueOf(name);
}

/// Optionally, enum_class can generate a mixin to go with your enum for use
/// with Angular. It exposes your enum constants as getters. So, if you mix it
/// in to your Dart component class, the values become available to the
/// corresponding Angular template.
///
/// Trigger mixin generation by writing a line like this one next to your enum.
abstract class RecognitionTaskStageMixin = Object with _$RecognitionTaskStageMixin;

