import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';

class AppStateProvider extends ChangeNotifier {
  final SharedPreferences prefs;

  late String _locale;
  late bool _isDarkMode;
  bool _notificationsEnabled = true;

  AppStateProvider(this.prefs) {
    _loadSettings();
  }

  // Getters
  String get locale => _locale;
  bool get isDarkMode => _isDarkMode;
  bool get notificationsEnabled => _notificationsEnabled;

  void _loadSettings() {
    _locale = prefs.getString('locale') ?? 'en';
    _isDarkMode = prefs.getBool('dark_mode') ?? false;
    _notificationsEnabled = prefs.getBool('notifications_enabled') ?? true;
  }

  // Set language/locale
  void setLocale(String languageCode) {
    _locale = languageCode;
    prefs.setString('locale', languageCode);
    notifyListeners();
  }

  // Toggle dark mode
  void toggleDarkMode() {
    _isDarkMode = !_isDarkMode;
    prefs.setBool('dark_mode', _isDarkMode);
    notifyListeners();
  }

  // Toggle notifications
  void toggleNotifications() {
    _notificationsEnabled = !_notificationsEnabled;
    prefs.setBool('notifications_enabled', _notificationsEnabled);
    notifyListeners();
  }

  // Mark onboarding as completed
  void completeOnboarding() {
    prefs.setBool('onboarding_completed', true);
    notifyListeners();
  }

  // Check if onboarding is completed
  bool isOnboardingCompleted() {
    return prefs.getBool('onboarding_completed') ?? false;
  }
}
