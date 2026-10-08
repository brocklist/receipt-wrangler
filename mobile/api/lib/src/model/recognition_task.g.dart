// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'recognition_task.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$RecognitionTask extends RecognitionTask {
  @override
  final int id;
  @override
  final String clientRequestId;
  @override
  final String fileName;
  @override
  final int fileSize;
  @override
  final int groupId;
  @override
  final int ownerUserId;
  @override
  final int version;
  @override
  final RecognitionTaskStatus status;
  @override
  final RecognitionTaskStage stage;
  @override
  final DateTime createdAt;
  @override
  final DateTime updatedAt;
  @override
  final DateTime? queuedAt;
  @override
  final DateTime? startedAt;
  @override
  final DateTime? stageStartedAt;
  @override
  final DateTime? completedAt;
  @override
  final int uploadedBytes;
  @override
  final int? uploadTotalBytes;
  @override
  final int attempt;
  @override
  final int maxAttempts;
  @override
  final DateTime? nextRetryAt;
  @override
  final bool fallbackActive;
  @override
  final int? receiptId;
  @override
  final String errorCode;
  @override
  final String errorMessage;
  @override
  final bool canUpload;
  @override
  final bool canRetry;

  factory _$RecognitionTask([void Function(RecognitionTaskBuilder)? updates]) =>
      (RecognitionTaskBuilder()..update(updates))._build();

  _$RecognitionTask._(
      {required this.id,
      required this.clientRequestId,
      required this.fileName,
      required this.fileSize,
      required this.groupId,
      required this.ownerUserId,
      required this.version,
      required this.status,
      required this.stage,
      required this.createdAt,
      required this.updatedAt,
      this.queuedAt,
      this.startedAt,
      this.stageStartedAt,
      this.completedAt,
      required this.uploadedBytes,
      this.uploadTotalBytes,
      required this.attempt,
      required this.maxAttempts,
      this.nextRetryAt,
      required this.fallbackActive,
      this.receiptId,
      required this.errorCode,
      required this.errorMessage,
      required this.canUpload,
      required this.canRetry})
      : super._();
  @override
  RecognitionTask rebuild(void Function(RecognitionTaskBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  RecognitionTaskBuilder toBuilder() => RecognitionTaskBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is RecognitionTask &&
        id == other.id &&
        clientRequestId == other.clientRequestId &&
        fileName == other.fileName &&
        fileSize == other.fileSize &&
        groupId == other.groupId &&
        ownerUserId == other.ownerUserId &&
        version == other.version &&
        status == other.status &&
        stage == other.stage &&
        createdAt == other.createdAt &&
        updatedAt == other.updatedAt &&
        queuedAt == other.queuedAt &&
        startedAt == other.startedAt &&
        stageStartedAt == other.stageStartedAt &&
        completedAt == other.completedAt &&
        uploadedBytes == other.uploadedBytes &&
        uploadTotalBytes == other.uploadTotalBytes &&
        attempt == other.attempt &&
        maxAttempts == other.maxAttempts &&
        nextRetryAt == other.nextRetryAt &&
        fallbackActive == other.fallbackActive &&
        receiptId == other.receiptId &&
        errorCode == other.errorCode &&
        errorMessage == other.errorMessage &&
        canUpload == other.canUpload &&
        canRetry == other.canRetry;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, clientRequestId.hashCode);
    _$hash = $jc(_$hash, fileName.hashCode);
    _$hash = $jc(_$hash, fileSize.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, ownerUserId.hashCode);
    _$hash = $jc(_$hash, version.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, stage.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, queuedAt.hashCode);
    _$hash = $jc(_$hash, startedAt.hashCode);
    _$hash = $jc(_$hash, stageStartedAt.hashCode);
    _$hash = $jc(_$hash, completedAt.hashCode);
    _$hash = $jc(_$hash, uploadedBytes.hashCode);
    _$hash = $jc(_$hash, uploadTotalBytes.hashCode);
    _$hash = $jc(_$hash, attempt.hashCode);
    _$hash = $jc(_$hash, maxAttempts.hashCode);
    _$hash = $jc(_$hash, nextRetryAt.hashCode);
    _$hash = $jc(_$hash, fallbackActive.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, errorCode.hashCode);
    _$hash = $jc(_$hash, errorMessage.hashCode);
    _$hash = $jc(_$hash, canUpload.hashCode);
    _$hash = $jc(_$hash, canRetry.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'RecognitionTask')
          ..add('id', id)
          ..add('clientRequestId', clientRequestId)
          ..add('fileName', fileName)
          ..add('fileSize', fileSize)
          ..add('groupId', groupId)
          ..add('ownerUserId', ownerUserId)
          ..add('version', version)
          ..add('status', status)
          ..add('stage', stage)
          ..add('createdAt', createdAt)
          ..add('updatedAt', updatedAt)
          ..add('queuedAt', queuedAt)
          ..add('startedAt', startedAt)
          ..add('stageStartedAt', stageStartedAt)
          ..add('completedAt', completedAt)
          ..add('uploadedBytes', uploadedBytes)
          ..add('uploadTotalBytes', uploadTotalBytes)
          ..add('attempt', attempt)
          ..add('maxAttempts', maxAttempts)
          ..add('nextRetryAt', nextRetryAt)
          ..add('fallbackActive', fallbackActive)
          ..add('receiptId', receiptId)
          ..add('errorCode', errorCode)
          ..add('errorMessage', errorMessage)
          ..add('canUpload', canUpload)
          ..add('canRetry', canRetry))
        .toString();
  }
}

class RecognitionTaskBuilder
    implements Builder<RecognitionTask, RecognitionTaskBuilder> {
  _$RecognitionTask? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _clientRequestId;
  String? get clientRequestId => _$this._clientRequestId;
  set clientRequestId(String? clientRequestId) =>
      _$this._clientRequestId = clientRequestId;

  String? _fileName;
  String? get fileName => _$this._fileName;
  set fileName(String? fileName) => _$this._fileName = fileName;

  int? _fileSize;
  int? get fileSize => _$this._fileSize;
  set fileSize(int? fileSize) => _$this._fileSize = fileSize;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _ownerUserId;
  int? get ownerUserId => _$this._ownerUserId;
  set ownerUserId(int? ownerUserId) => _$this._ownerUserId = ownerUserId;

  int? _version;
  int? get version => _$this._version;
  set version(int? version) => _$this._version = version;

  RecognitionTaskStatus? _status;
  RecognitionTaskStatus? get status => _$this._status;
  set status(RecognitionTaskStatus? status) => _$this._status = status;

  RecognitionTaskStage? _stage;
  RecognitionTaskStage? get stage => _$this._stage;
  set stage(RecognitionTaskStage? stage) => _$this._stage = stage;

  DateTime? _createdAt;
  DateTime? get createdAt => _$this._createdAt;
  set createdAt(DateTime? createdAt) => _$this._createdAt = createdAt;

  DateTime? _updatedAt;
  DateTime? get updatedAt => _$this._updatedAt;
  set updatedAt(DateTime? updatedAt) => _$this._updatedAt = updatedAt;

  DateTime? _queuedAt;
  DateTime? get queuedAt => _$this._queuedAt;
  set queuedAt(DateTime? queuedAt) => _$this._queuedAt = queuedAt;

  DateTime? _startedAt;
  DateTime? get startedAt => _$this._startedAt;
  set startedAt(DateTime? startedAt) => _$this._startedAt = startedAt;

  DateTime? _stageStartedAt;
  DateTime? get stageStartedAt => _$this._stageStartedAt;
  set stageStartedAt(DateTime? stageStartedAt) =>
      _$this._stageStartedAt = stageStartedAt;

  DateTime? _completedAt;
  DateTime? get completedAt => _$this._completedAt;
  set completedAt(DateTime? completedAt) => _$this._completedAt = completedAt;

  int? _uploadedBytes;
  int? get uploadedBytes => _$this._uploadedBytes;
  set uploadedBytes(int? uploadedBytes) =>
      _$this._uploadedBytes = uploadedBytes;

  int? _uploadTotalBytes;
  int? get uploadTotalBytes => _$this._uploadTotalBytes;
  set uploadTotalBytes(int? uploadTotalBytes) =>
      _$this._uploadTotalBytes = uploadTotalBytes;

  int? _attempt;
  int? get attempt => _$this._attempt;
  set attempt(int? attempt) => _$this._attempt = attempt;

  int? _maxAttempts;
  int? get maxAttempts => _$this._maxAttempts;
  set maxAttempts(int? maxAttempts) => _$this._maxAttempts = maxAttempts;

  DateTime? _nextRetryAt;
  DateTime? get nextRetryAt => _$this._nextRetryAt;
  set nextRetryAt(DateTime? nextRetryAt) => _$this._nextRetryAt = nextRetryAt;

  bool? _fallbackActive;
  bool? get fallbackActive => _$this._fallbackActive;
  set fallbackActive(bool? fallbackActive) =>
      _$this._fallbackActive = fallbackActive;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  String? _errorCode;
  String? get errorCode => _$this._errorCode;
  set errorCode(String? errorCode) => _$this._errorCode = errorCode;

  String? _errorMessage;
  String? get errorMessage => _$this._errorMessage;
  set errorMessage(String? errorMessage) => _$this._errorMessage = errorMessage;

  bool? _canUpload;
  bool? get canUpload => _$this._canUpload;
  set canUpload(bool? canUpload) => _$this._canUpload = canUpload;

  bool? _canRetry;
  bool? get canRetry => _$this._canRetry;
  set canRetry(bool? canRetry) => _$this._canRetry = canRetry;

  RecognitionTaskBuilder() {
    RecognitionTask._defaults(this);
  }

  RecognitionTaskBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _clientRequestId = $v.clientRequestId;
      _fileName = $v.fileName;
      _fileSize = $v.fileSize;
      _groupId = $v.groupId;
      _ownerUserId = $v.ownerUserId;
      _version = $v.version;
      _status = $v.status;
      _stage = $v.stage;
      _createdAt = $v.createdAt;
      _updatedAt = $v.updatedAt;
      _queuedAt = $v.queuedAt;
      _startedAt = $v.startedAt;
      _stageStartedAt = $v.stageStartedAt;
      _completedAt = $v.completedAt;
      _uploadedBytes = $v.uploadedBytes;
      _uploadTotalBytes = $v.uploadTotalBytes;
      _attempt = $v.attempt;
      _maxAttempts = $v.maxAttempts;
      _nextRetryAt = $v.nextRetryAt;
      _fallbackActive = $v.fallbackActive;
      _receiptId = $v.receiptId;
      _errorCode = $v.errorCode;
      _errorMessage = $v.errorMessage;
      _canUpload = $v.canUpload;
      _canRetry = $v.canRetry;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(RecognitionTask other) {
    _$v = other as _$RecognitionTask;
  }

  @override
  void update(void Function(RecognitionTaskBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  RecognitionTask build() => _build();

  _$RecognitionTask _build() {
    final _$result = _$v ??
        _$RecognitionTask._(
          id: BuiltValueNullFieldError.checkNotNull(
              id, r'RecognitionTask', 'id'),
          clientRequestId: BuiltValueNullFieldError.checkNotNull(
              clientRequestId, r'RecognitionTask', 'clientRequestId'),
          fileName: BuiltValueNullFieldError.checkNotNull(
              fileName, r'RecognitionTask', 'fileName'),
          fileSize: BuiltValueNullFieldError.checkNotNull(
              fileSize, r'RecognitionTask', 'fileSize'),
          groupId: BuiltValueNullFieldError.checkNotNull(
              groupId, r'RecognitionTask', 'groupId'),
          ownerUserId: BuiltValueNullFieldError.checkNotNull(
              ownerUserId, r'RecognitionTask', 'ownerUserId'),
          version: BuiltValueNullFieldError.checkNotNull(
              version, r'RecognitionTask', 'version'),
          status: BuiltValueNullFieldError.checkNotNull(
              status, r'RecognitionTask', 'status'),
          stage: BuiltValueNullFieldError.checkNotNull(
              stage, r'RecognitionTask', 'stage'),
          createdAt: BuiltValueNullFieldError.checkNotNull(
              createdAt, r'RecognitionTask', 'createdAt'),
          updatedAt: BuiltValueNullFieldError.checkNotNull(
              updatedAt, r'RecognitionTask', 'updatedAt'),
          queuedAt: queuedAt,
          startedAt: startedAt,
          stageStartedAt: stageStartedAt,
          completedAt: completedAt,
          uploadedBytes: BuiltValueNullFieldError.checkNotNull(
              uploadedBytes, r'RecognitionTask', 'uploadedBytes'),
          uploadTotalBytes: uploadTotalBytes,
          attempt: BuiltValueNullFieldError.checkNotNull(
              attempt, r'RecognitionTask', 'attempt'),
          maxAttempts: BuiltValueNullFieldError.checkNotNull(
              maxAttempts, r'RecognitionTask', 'maxAttempts'),
          nextRetryAt: nextRetryAt,
          fallbackActive: BuiltValueNullFieldError.checkNotNull(
              fallbackActive, r'RecognitionTask', 'fallbackActive'),
          receiptId: receiptId,
          errorCode: BuiltValueNullFieldError.checkNotNull(
              errorCode, r'RecognitionTask', 'errorCode'),
          errorMessage: BuiltValueNullFieldError.checkNotNull(
              errorMessage, r'RecognitionTask', 'errorMessage'),
          canUpload: BuiltValueNullFieldError.checkNotNull(
              canUpload, r'RecognitionTask', 'canUpload'),
          canRetry: BuiltValueNullFieldError.checkNotNull(
              canRetry, r'RecognitionTask', 'canRetry'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
