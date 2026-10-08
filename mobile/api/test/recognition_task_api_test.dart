import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for RecognitionTaskApi
void main() {
  final instance = Openapi().getRecognitionTaskApi();

  group(RecognitionTaskApi, () {
    // Register a Quick Scan file before uploading
    //
    //Future<RecognitionTask> createRecognitionTask(CreateRecognitionTaskCommand createRecognitionTaskCommand) async
    test('test createRecognitionTask', () async {
      // TODO
    });

    // Refresh one authorized task
    //
    //Future<RecognitionTask> getRecognitionTask(int id) async
    test('test getRecognitionTask', () async {
      // TODO
    });

    // Get authorized Quick Scan tasks and scoped counts
    //
    //Future<GetRecognitionTasksResponse> getRecognitionTasks({ String scope, String bucket, int page, int pageSize, String clientRequestId, String ids, int groupId, RecognitionTaskStatus status }) async
    test('test getRecognitionTasks', () async {
      // TODO
    });

    // Retry a failed task using its retained source
    //
    //Future<RecognitionTask> retryRecognitionTask(int id, RetryRecognitionTaskCommand retryRecognitionTaskCommand) async
    test('test retryRecognitionTask', () async {
      // TODO
    });

    // Upload a complete Quick Scan source file
    //
    //Future<RecognitionTask> uploadRecognitionTaskFile(int id, MultipartFile file) async
    test('test uploadRecognitionTaskFile', () async {
      // TODO
    });

  });
}
