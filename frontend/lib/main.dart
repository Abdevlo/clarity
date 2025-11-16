import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'config/localization.dart';
import 'config/theme.dart';
import 'providers/auth_provider.dart';
import 'providers/app_state_provider.dart';
import 'providers/health_provider.dart';
import 'screens/onboarding_screen.dart';
import 'screens/auth/login_screen.dart';
import 'screens/home_screen.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  final prefs = await SharedPreferences.getInstance();

  runApp(MyApp(prefs: prefs));
}

class MyApp extends StatefulWidget {
  final SharedPreferences prefs;

  const MyApp({required this.prefs});

  @override
  State<MyApp> createState() => _MyAppState();
}

class _MyAppState extends State<MyApp> {
  late AppStateProvider appStateProvider;
  late AuthProvider authProvider;
  late HealthProvider healthProvider;

  @override
  void initState() {
    super.initState();

    appStateProvider = AppStateProvider(widget.prefs);
    authProvider = AuthProvider(widget.prefs);
    healthProvider = HealthProvider();

    // Check if user is authenticated
    _checkAuthStatus();
  }

  void _checkAuthStatus() async {
    final token = widget.prefs.getString('access_token');
    await Future.delayed(Duration(seconds: 1)); // Simulate splash screen

    if (token != null) {
      authProvider.setToken(token);
    }
  }

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: appStateProvider),
        ChangeNotifierProvider.value(value: authProvider),
        ChangeNotifierProvider.value(value: healthProvider),
      ],
      child: MaterialApp(
        title: 'Clarity',
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
        themeMode: appStateProvider.isDarkMode ? ThemeMode.dark : ThemeMode.light,

        localizationsDelegates: [
          AppLocalizations.delegate,
          GlobalMaterialLocalizations.delegate,
          GlobalWidgetsLocalizations.delegate,
          GlobalCupertinoLocalizations.delegate,
        ],
        supportedLocales: [
          Locale('en', ''),
          Locale('es', ''),
        ],
        locale: Locale(appStateProvider.locale),

        home: _buildHome(),
      ),
    );
  }

  Widget _buildHome() {
    return Consumer<AuthProvider>(
      builder: (context, auth, _) {
        // Check onboarding status
        if (!widget.prefs.containsKey('onboarding_completed')) {
          return OnboardingScreen();
        }

        // Check authentication status
        if (auth.isAuthenticated) {
          return HomeScreen();
        }

        return LoginScreen();
      },
    );
  }
}
