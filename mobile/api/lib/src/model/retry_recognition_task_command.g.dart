// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'retry_recognition_task_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$RetryRecognitionTaskCommand extends RetryRecognitionTaskCommand {
  @override
  final int version;

  factory _$RetryRecognitionTaskCommand(
          [void Function(RetryRecognitionTaskCommandBuilder)? updates]) =>
      (RetryRecognitionTaskCommandBuilder()..update(updates))._build();

  _$RetryRecognitionTaskCommand._({required this.version}) : super._();
  @override
  RetryRecognitionTaskCommand rebuild(
          void Function(RetryRecognitionTaskCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  RetryRecognitionTaskCommandBuilder toBuilder() =>
      RetryRecognitionTaskCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is RetryRecognitionTaskCommand && version == other.version;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, version.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'RetryRecognitionTaskCommand')
          ..add('version', version))
        .toString();
  }
}

class RetryRecognitionTaskCommandBuilder
    implements
        Builder<RetryRecognitionTaskCommand,
            RetryRecognitionTaskCommandBuilder> {
  _$RetryRecognitionTaskCommand? _$v;

  int? _version;
  int? get version => _$this._version;
  set version(int? version) => _$this._version = version;

  RetryRecognitionTaskCommandBuilder() {
    RetryRecognitionTaskCommand._defaults(this);
  }

  RetryRecognitionTaskCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _version = $v.version;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(RetryRecognitionTaskCommand other) {
    _$v = other as _$RetryRecognitionTaskCommand;
  }

  @override
  void update(void Function(RetryRecognitionTaskCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  RetryRecognitionTaskCommand build() => _build();

  _$RetryRecognitionTaskCommand _build() {
    final _$result = _$v ??
        _$RetryRecognitionTaskCommand._(
          version: BuiltValueNullFieldError.checkNotNull(
              version, r'RetryRecognitionTaskCommand', 'version'),
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
