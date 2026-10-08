// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'get_recognition_tasks_response.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$GetRecognitionTasksResponse extends GetRecognitionTasksResponse {
  @override
  final BuiltList<RecognitionTask> data;
  @override
  final int totalCount;
  @override
  final int activeCount;
  @override
  final int awaitingUploadCount;
  @override
  final int runningCount;
  @override
  final int failedCount;

  factory _$GetRecognitionTasksResponse(
          [void Function(GetRecognitionTasksResponseBuilder)? updates]) =>
      (GetRecognitionTasksResponseBuilder()..update(updates))._build();

  _$GetRecognitionTasksResponse._(
      {required this.data,
      required this.totalCount,
      required this.activeCount,
      required this.awaitingUploadCount,
      required this.runningCount,
      required this.failedCount})
      : super._();
  @override
  GetRecognitionTasksResponse rebuild(
          void Function(GetRecognitionTasksResponseBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GetRecognitionTasksResponseBuilder toBuilder() =>
      GetRecognitionTasksResponseBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GetRecognitionTasksResponse &&
        data == other.data &&
        totalCount == other.totalCount &&
        activeCount == other.activeCount &&
        awaitingUploadCount == other.awaitingUploadCount &&
        runningCount == other.runningCount &&
        failedCount == other.failedCount;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, data.hashCode);
    _$hash = $jc(_$hash, totalCount.hashCode);
    _$hash = $jc(_$hash, activeCount.hashCode);
    _$hash = $jc(_$hash, awaitingUploadCount.hashCode);
    _$hash = $jc(_$hash, runningCount.hashCode);
    _$hash = $jc(_$hash, failedCount.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'GetRecognitionTasksResponse')
          ..add('data', data)
          ..add('totalCount', totalCount)
          ..add('activeCount', activeCount)
          ..add('awaitingUploadCount', awaitingUploadCount)
          ..add('runningCount', runningCount)
          ..add('failedCount', failedCount))
        .toString();
  }
}

class GetRecognitionTasksResponseBuilder
    implements
        Builder<GetRecognitionTasksResponse,
            GetRecognitionTasksResponseBuilder> {
  _$GetRecognitionTasksResponse? _$v;

  ListBuilder<RecognitionTask>? _data;
  ListBuilder<RecognitionTask> get data =>
      _$this._data ??= ListBuilder<RecognitionTask>();
  set data(ListBuilder<RecognitionTask>? data) => _$this._data = data;

  int? _totalCount;
  int? get totalCount => _$this._totalCount;
  set totalCount(int? totalCount) => _$this._totalCount = totalCount;

  int? _activeCount;
  int? get activeCount => _$this._activeCount;
  set activeCount(int? activeCount) => _$this._activeCount = activeCount;

  int? _awaitingUploadCount;
  int? get awaitingUploadCount => _$this._awaitingUploadCount;
  set awaitingUploadCount(int? awaitingUploadCount) =>
      _$this._awaitingUploadCount = awaitingUploadCount;

  int? _runningCount;
  int? get runningCount => _$this._runningCount;
  set runningCount(int? runningCount) => _$this._runningCount = runningCount;

  int? _failedCount;
  int? get failedCount => _$this._failedCount;
  set failedCount(int? failedCount) => _$this._failedCount = failedCount;

  GetRecognitionTasksResponseBuilder() {
    GetRecognitionTasksResponse._defaults(this);
  }

  GetRecognitionTasksResponseBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _data = $v.data.toBuilder();
      _totalCount = $v.totalCount;
      _activeCount = $v.activeCount;
      _awaitingUploadCount = $v.awaitingUploadCount;
      _runningCount = $v.runningCount;
      _failedCount = $v.failedCount;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(GetRecognitionTasksResponse other) {
    _$v = other as _$GetRecognitionTasksResponse;
  }

  @override
  void update(void Function(GetRecognitionTasksResponseBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GetRecognitionTasksResponse build() => _build();

  _$GetRecognitionTasksResponse _build() {
    _$GetRecognitionTasksResponse _$result;
    try {
      _$result = _$v ??
          _$GetRecognitionTasksResponse._(
            data: data.build(),
            totalCount: BuiltValueNullFieldError.checkNotNull(
                totalCount, r'GetRecognitionTasksResponse', 'totalCount'),
            activeCount: BuiltValueNullFieldError.checkNotNull(
                activeCount, r'GetRecognitionTasksResponse', 'activeCount'),
            awaitingUploadCount: BuiltValueNullFieldError.checkNotNull(
                awaitingUploadCount,
                r'GetRecognitionTasksResponse',
                'awaitingUploadCount'),
            runningCount: BuiltValueNullFieldError.checkNotNull(
                runningCount, r'GetRecognitionTasksResponse', 'runningCount'),
            failedCount: BuiltValueNullFieldError.checkNotNull(
                failedCount, r'GetRecognitionTasksResponse', 'failedCount'),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'data';
        data.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'GetRecognitionTasksResponse', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
