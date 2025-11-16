import 'package:flutter/foundation.dart';

class HealthRecord {
  final String id;
  final String userId;
  final String recordType;
  final String title;
  final String description;
  final Map<String, String> metadata;
  final DateTime createdAt;
  final DateTime updatedAt;

  HealthRecord({
    required this.id,
    required this.userId,
    required this.recordType,
    required this.title,
    required this.description,
    required this.metadata,
    required this.createdAt,
    required this.updatedAt,
  });
}

class HealthProvider extends ChangeNotifier {
  List<HealthRecord> _records = [];
  bool _isLoading = false;
  String? _error;

  // Getters
  List<HealthRecord> get records => _records;
  bool get isLoading => _isLoading;
  String? get error => _error;

  // Fetch records
  Future<void> fetchRecords(String userId, {int limit = 20, int offset = 0}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      // TODO: Implement gRPC call to backend
      // Mock data
      await Future.delayed(Duration(seconds: 1));

      _records = [
        HealthRecord(
          id: '1',
          userId: userId,
          recordType: 'prescription',
          title: 'Aspirin Prescription',
          description: '500mg tablets',
          metadata: {'dosage': '500mg', 'frequency': 'Twice daily'},
          createdAt: DateTime.now().subtract(Duration(days: 2)),
          updatedAt: DateTime.now().subtract(Duration(days: 2)),
        ),
      ];

      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
    }
  }

  // Create record
  Future<bool> createRecord({
    required String userId,
    required String recordType,
    required String title,
    required String description,
    required Map<String, String> metadata,
  }) async {
    _isLoading = true;
    notifyListeners();

    try {
      // TODO: Implement gRPC call to backend
      await Future.delayed(Duration(seconds: 1));

      _records.insert(
        0,
        HealthRecord(
          id: DateTime.now().millisecondsSinceEpoch.toString(),
          userId: userId,
          recordType: recordType,
          title: title,
          description: description,
          metadata: metadata,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
      );

      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  // Delete record
  Future<bool> deleteRecord(String recordId) async {
    try {
      // TODO: Implement gRPC call to backend
      _records.removeWhere((r) => r.id == recordId);
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  // Clear error
  void clearError() {
    _error = null;
    notifyListeners();
  }
}
