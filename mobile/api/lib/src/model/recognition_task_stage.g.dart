// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'recognition_task_stage.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const RecognitionTaskStage _$UPLOAD = const RecognitionTaskStage._('UPLOAD');
const RecognitionTaskStage _$PREPROCESSING =
    const RecognitionTaskStage._('PREPROCESSING');
const RecognitionTaskStage _$OCR = const RecognitionTaskStage._('OCR');
const RecognitionTaskStage _$AI = const RecognitionTaskStage._('AI');
const RecognitionTaskStage _$PARSING = const RecognitionTaskStage._('PARSING');
const RecognitionTaskStage _$SAVING = const RecognitionTaskStage._('SAVING');
const RecognitionTaskStage _$DONE = const RecognitionTaskStage._('DONE');

RecognitionTaskStage _$valueOf(String name) {
  switch (name) {
    case 'UPLOAD':
      return _$UPLOAD;
    case 'PREPROCESSING':
      return _$PREPROCESSING;
    case 'OCR':
      return _$OCR;
    case 'AI':
      return _$AI;
    case 'PARSING':
      return _$PARSING;
    case 'SAVING':
      return _$SAVING;
    case 'DONE':
      return _$DONE;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<RecognitionTaskStage> _$values =
    BuiltSet<RecognitionTaskStage>(const <RecognitionTaskStage>[
  _$UPLOAD,
  _$PREPROCESSING,
  _$OCR,
  _$AI,
  _$PARSING,
  _$SAVING,
  _$DONE,
]);

class _$RecognitionTaskStageMeta {
  const _$RecognitionTaskStageMeta();
  RecognitionTaskStage get UPLOAD => _$UPLOAD;
  RecognitionTaskStage get PREPROCESSING => _$PREPROCESSING;
  RecognitionTaskStage get OCR => _$OCR;
  RecognitionTaskStage get AI => _$AI;
  RecognitionTaskStage get PARSING => _$PARSING;
  RecognitionTaskStage get SAVING => _$SAVING;
  RecognitionTaskStage get DONE => _$DONE;
  RecognitionTaskStage valueOf(String name) => _$valueOf(name);
  BuiltSet<RecognitionTaskStage> get values => _$values;
}

abstract class _$RecognitionTaskStageMixin {
  // ignore: non_constant_identifier_names
  _$RecognitionTaskStageMeta get RecognitionTaskStage =>
      const _$RecognitionTaskStageMeta();
}

Serializer<RecognitionTaskStage> _$recognitionTaskStageSerializer =
    _$RecognitionTaskStageSerializer();

class _$RecognitionTaskStageSerializer
    implements PrimitiveSerializer<RecognitionTaskStage> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'UPLOAD': 'UPLOAD',
    'PREPROCESSING': 'PREPROCESSING',
    'OCR': 'OCR',
    'AI': 'AI',
    'PARSING': 'PARSING',
    'SAVING': 'SAVING',
    'DONE': 'DONE',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'UPLOAD': 'UPLOAD',
    'PREPROCESSING': 'PREPROCESSING',
    'OCR': 'OCR',
    'AI': 'AI',
    'PARSING': 'PARSING',
    'SAVING': 'SAVING',
    'DONE': 'DONE',
  };

  @override
  final Iterable<Type> types = const <Type>[RecognitionTaskStage];
  @override
  final String wireName = 'RecognitionTaskStage';

  @override
  Object serialize(Serializers serializers, RecognitionTaskStage object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  RecognitionTaskStage deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      RecognitionTaskStage.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
