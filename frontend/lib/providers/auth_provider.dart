import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

enum UserRole { patient, doctor }

class AuthProvider extends ChangeNotifier {
  final SharedPreferences prefs;

  String? _accessToken;
  String? _refreshToken;
  String? _userId;
  String? _userEmail;
  bool _isDoctorMode = false;
  bool _isLoading = false;
  String? _error;

  AuthProvider(this.prefs) {
    _loadFromPrefs();
  }

  // Getters
  String? get accessToken => _accessToken;
  String? get refreshToken => _refreshToken;
  String? get userId => _userId;
  String? get userEmail => _userEmail;
  bool get isDoctorMode => _isDoctorMode;
  bool get isAuthenticated => _accessToken != null;
  bool get isLoading => _isLoading;
  String? get error => _error;

  // Load from SharedPreferences
  void _loadFromPrefs() {
    _accessToken = prefs.getString('access_token');
    _refreshToken = prefs.getString('refresh_token');
    _userId = prefs.getString('user_id');
    _userEmail = prefs.getString('user_email');
    _isDoctorMode = prefs.getBool('doctor_mode') ?? false;
  }

  // Set token
  void setToken(String token) {
    _accessToken = token;
    prefs.setString('access_token', token);
    notifyListeners();
  }

  // Set user data
  void setUser(String id, String email, {String? name}) {
    _userId = id;
    _userEmail = email;
    prefs.setString('user_id', id);
    prefs.setString('user_email', email);
    if (name != null) {
      prefs.setString('user_name', name);
    }
    notifyListeners();
  }

  // Send OTP
  Future<bool> sendOTP(String email) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      // TODO: Implement gRPC call to backend
      // Mock delay
      await Future.delayed(Duration(seconds: 2));

      _error = null;
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

  // Verify OTP
  Future<bool> verifyOTP(String email, String otp) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      // TODO: Implement gRPC call to backend
      // Mock response
      _accessToken = 'mock_token_${DateTime.now().millisecondsSinceEpoch}';
      _refreshToken = 'mock_refresh_token_${DateTime.now().millisecondsSinceEpoch}';
      _userId = 'user_123';
      _userEmail = email;

      // Save to preferences
      prefs.setString('access_token', _accessToken!);
      prefs.setString('refresh_token', _refreshToken!);
      prefs.setString('user_id', _userId!);
      prefs.setString('user_email', _userEmail!);

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

  // Logout
  void logout() {
    _accessToken = null;
    _refreshToken = null;
    _userId = null;
    _userEmail = null;

    prefs.remove('access_token');
    prefs.remove('refresh_token');
    prefs.remove('user_id');
    prefs.remove('user_email');

    notifyListeners();
  }

  // Clear error
  void clearError() {
    _error = null;
    notifyListeners();
  }

  // Toggle doctor mode
  void setDoctorMode(bool value) {
    _isDoctorMode = value;
    prefs.setBool('doctor_mode', value);
    notifyListeners();
  }
}
