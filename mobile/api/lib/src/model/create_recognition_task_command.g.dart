// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'create_recognition_task_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CreateRecognitionTaskCommand extends CreateRecognitionTaskCommand {
  @override
  final String clientRequestId;
  @override
  final String fileName;
  @override
  final int fileSize;
  @override
  final int groupId;
  @override
  final int? paidByUserId;
  @override
  final ReceiptStatus? status;
  @override
  final BuiltList<int>? categoryIds;
  @override
  final BuiltList<int>? tagIds;
  @override
  final String? comment;

  factory _$CreateRecognitionTaskCommand(
          [void Function(CreateRecognitionTaskCommandBuilder)? updates]) =>
      (CreateRecognitionTaskCommandBuilder()..update(updates))._build();

  _$CreateRecognitionTaskCommand._(
      {required this.clientRequestId,
      required this.fileName,
      required this.fileSize,
      required this.groupId,
      this.paidByUserId,
      this.status,
      this.categoryIds,
      this.tagIds,
      this.comment})
      : super._();
  @override
  CreateRecognitionTaskCommand rebuild(
          void Function(CreateRecognitionTaskCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CreateRecognitionTaskCommandBuilder toBuilder() =>
      CreateRecognitionTaskCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CreateRecognitionTaskCommand &&
        clientRequestId == other.clientRequestId &&
        fileName == other.fileName &&
        fileSize == other.fileSize &&
        groupId == other.groupId &&
        paidByUserId == other.paidByUserId &&
        status == other.status &&
        categoryIds == other.categoryIds &&
        tagIds == other.tagIds &&
        comment == other.comment;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, clientRequestId.hashCode);
    _$hash = $jc(_$hash, fileName.hashCode);
    _$hash = $jc(_$hash, fileSize.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, paidByUserId.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, categoryIds.hashCode);
    _$hash = $jc(_$hash, tagIds.hashCode);
    _$hash = $jc(_$hash, comment.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'CreateRecognitionTaskCommand')
          ..add('clientRequestId', clientRequestId)
          ..add('fileName', fileName)
          ..add('fileSize', fileSize)
          ..add('groupId', groupId)
          ..add('paidByUserId', paidByUserId)
          ..add('status', status)
          ..add('categoryIds', categoryIds)
          ..add('tagIds', tagIds)
          ..add('comment', comment))
        .toString();
  }
}

class CreateRecognitionTaskCommandBuilder
    implements
        Builder<CreateRecognitionTaskCommand,
            CreateRecognitionTaskCommandBuilder> {
  _$CreateRecognitionTaskCommand? _$v;

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

  int? _paidByUserId;
  int? get paidByUserId => _$this._paidByUserId;
  set paidByUserId(int? paidByUserId) => _$this._paidByUserId = paidByUserId;

  ReceiptStatus? _status;
  ReceiptStatus? get status => _$this._status;
  set status(ReceiptStatus? status) => _$this._status = status;

  ListBuilder<int>? _categoryIds;
  ListBuilder<int> get categoryIds =>
      _$this._categoryIds ??= ListBuilder<int>();
  set categoryIds(ListBuilder<int>? categoryIds) =>
      _$this._categoryIds = categoryIds;

  ListBuilder<int>? _tagIds;
  ListBuilder<int> get tagIds => _$this._tagIds ??= ListBuilder<int>();
  set tagIds(ListBuilder<int>? tagIds) => _$this._tagIds = tagIds;

  String? _comment;
  String? get comment => _$this._comment;
  set comment(String? comment) => _$this._comment = comment;

  CreateRecognitionTaskCommandBuilder() {
    CreateRecognitionTaskCommand._defaults(this);
  }

  CreateRecognitionTaskCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _clientRequestId = $v.clientRequestId;
      _fileName = $v.fileName;
      _fileSize = $v.fileSize;
      _groupId = $v.groupId;
      _paidByUserId = $v.paidByUserId;
      _status = $v.status;
      _categoryIds = $v.categoryIds?.toBuilder();
      _tagIds = $v.tagIds?.toBuilder();
      _comment = $v.comment;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(CreateRecognitionTaskCommand other) {
    _$v = other as _$CreateRecognitionTaskCommand;
  }

  @override
  void update(void Function(CreateRecognitionTaskCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CreateRecognitionTaskCommand build() => _build();

  _$CreateRecognitionTaskCommand _build() {
    _$CreateRecognitionTaskCommand _$result;
    try {
      _$result = _$v ??
          _$CreateRecognitionTaskCommand._(
            clientRequestId: BuiltValueNullFieldError.checkNotNull(
                clientRequestId,
                r'CreateRecognitionTaskCommand',
                'clientRequestId'),
            fileName: BuiltValueNullFieldError.checkNotNull(
                fileName, r'CreateRecognitionTaskCommand', 'fileName'),
            fileSize: BuiltValueNullFieldError.checkNotNull(
                fileSize, r'CreateRecognitionTaskCommand', 'fileSize'),
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'CreateRecognitionTaskCommand', 'groupId'),
            paidByUserId: paidByUserId,
            status: status,
            categoryIds: _categoryIds?.build(),
            tagIds: _tagIds?.build(),
            comment: comment,
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'categoryIds';
        _categoryIds?.build();
        _$failedField = 'tagIds';
        _tagIds?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'CreateRecognitionTaskCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
