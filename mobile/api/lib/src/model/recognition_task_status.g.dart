// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'recognition_task_status.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

const RecognitionTaskStatus _$AWAITING_UPLOAD =
    const RecognitionTaskStatus._('AWAITING_UPLOAD');
const RecognitionTaskStatus _$UPLOADING =
    const RecognitionTaskStatus._('UPLOADING');
const RecognitionTaskStatus _$UPLOAD_INTERRUPTED =
    const RecognitionTaskStatus._('UPLOAD_INTERRUPTED');
const RecognitionTaskStatus _$DISPATCH_PENDING =
    const RecognitionTaskStatus._('DISPATCH_PENDING');
const RecognitionTaskStatus _$QUEUED = const RecognitionTaskStatus._('QUEUED');
const RecognitionTaskStatus _$RUNNING =
    const RecognitionTaskStatus._('RUNNING');
const RecognitionTaskStatus _$RETRY_WAIT =
    const RecognitionTaskStatus._('RETRY_WAIT');
const RecognitionTaskStatus _$SUCCEEDED =
    const RecognitionTaskStatus._('SUCCEEDED');
const RecognitionTaskStatus _$FAILED = const RecognitionTaskStatus._('FAILED');

RecognitionTaskStatus _$valueOf(String name) {
  switch (name) {
    case 'AWAITING_UPLOAD':
      return _$AWAITING_UPLOAD;
    case 'UPLOADING':
      return _$UPLOADING;
    case 'UPLOAD_INTERRUPTED':
      return _$UPLOAD_INTERRUPTED;
    case 'DISPATCH_PENDING':
      return _$DISPATCH_PENDING;
    case 'QUEUED':
      return _$QUEUED;
    case 'RUNNING':
      return _$RUNNING;
    case 'RETRY_WAIT':
      return _$RETRY_WAIT;
    case 'SUCCEEDED':
      return _$SUCCEEDED;
    case 'FAILED':
      return _$FAILED;
    default:
      throw ArgumentError(name);
  }
}

final BuiltSet<RecognitionTaskStatus> _$values =
    BuiltSet<RecognitionTaskStatus>(const <RecognitionTaskStatus>[
  _$AWAITING_UPLOAD,
  _$UPLOADING,
  _$UPLOAD_INTERRUPTED,
  _$DISPATCH_PENDING,
  _$QUEUED,
  _$RUNNING,
  _$RETRY_WAIT,
  _$SUCCEEDED,
  _$FAILED,
]);

class _$RecognitionTaskStatusMeta {
  const _$RecognitionTaskStatusMeta();
  RecognitionTaskStatus get AWAITING_UPLOAD => _$AWAITING_UPLOAD;
  RecognitionTaskStatus get UPLOADING => _$UPLOADING;
  RecognitionTaskStatus get UPLOAD_INTERRUPTED => _$UPLOAD_INTERRUPTED;
  RecognitionTaskStatus get DISPATCH_PENDING => _$DISPATCH_PENDING;
  RecognitionTaskStatus get QUEUED => _$QUEUED;
  RecognitionTaskStatus get RUNNING => _$RUNNING;
  RecognitionTaskStatus get RETRY_WAIT => _$RETRY_WAIT;
  RecognitionTaskStatus get SUCCEEDED => _$SUCCEEDED;
  RecognitionTaskStatus get FAILED => _$FAILED;
  RecognitionTaskStatus valueOf(String name) => _$valueOf(name);
  BuiltSet<RecognitionTaskStatus> get values => _$values;
}

abstract class _$RecognitionTaskStatusMixin {
  // ignore: non_constant_identifier_names
  _$RecognitionTaskStatusMeta get RecognitionTaskStatus =>
      const _$RecognitionTaskStatusMeta();
}

Serializer<RecognitionTaskStatus> _$recognitionTaskStatusSerializer =
    _$RecognitionTaskStatusSerializer();

class _$RecognitionTaskStatusSerializer
    implements PrimitiveSerializer<RecognitionTaskStatus> {
  static const Map<String, Object> _toWire = const <String, Object>{
    'AWAITING_UPLOAD': 'AWAITING_UPLOAD',
    'UPLOADING': 'UPLOADING',
    'UPLOAD_INTERRUPTED': 'UPLOAD_INTERRUPTED',
    'DISPATCH_PENDING': 'DISPATCH_PENDING',
    'QUEUED': 'QUEUED',
    'RUNNING': 'RUNNING',
    'RETRY_WAIT': 'RETRY_WAIT',
    'SUCCEEDED': 'SUCCEEDED',
    'FAILED': 'FAILED',
  };
  static const Map<Object, String> _fromWire = const <Object, String>{
    'AWAITING_UPLOAD': 'AWAITING_UPLOAD',
    'UPLOADING': 'UPLOADING',
    'UPLOAD_INTERRUPTED': 'UPLOAD_INTERRUPTED',
    'DISPATCH_PENDING': 'DISPATCH_PENDING',
    'QUEUED': 'QUEUED',
    'RUNNING': 'RUNNING',
    'RETRY_WAIT': 'RETRY_WAIT',
    'SUCCEEDED': 'SUCCEEDED',
    'FAILED': 'FAILED',
  };

  @override
  final Iterable<Type> types = const <Type>[RecognitionTaskStatus];
  @override
  final String wireName = 'RecognitionTaskStatus';

  @override
  Object serialize(Serializers serializers, RecognitionTaskStatus object,
          {FullType specifiedType = FullType.unspecified}) =>
      _toWire[object.name] ?? object.name;

  @override
  RecognitionTaskStatus deserialize(Serializers serializers, Object serialized,
          {FullType specifiedType = FullType.unspecified}) =>
      RecognitionTaskStatus.valueOf(
          _fromWire[serialized] ?? (serialized is String ? serialized : ''));
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
