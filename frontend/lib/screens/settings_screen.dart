import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../config/localization.dart';
import '../config/theme.dart';
import '../providers/app_state_provider.dart';
import '../providers/auth_provider.dart';
import 'auth/login_screen.dart';

class SettingsScreen extends StatefulWidget {
  @override
  State<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends State<SettingsScreen> {
  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;
    final appState = Provider.of<AppStateProvider>(context);
    final authProvider = Provider.of<AuthProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(context.translate('settings_title')),
        elevation: 0,
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: EdgeInsets.all(isMobile ? 12 : 24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Language Settings
              Card(
                child: Column(
                  children: [
                    ListTile(
                      leading: Icon(Icons.language),
                      title: Text(context.translate('settings_language')),
                      trailing: DropdownButton<String>(
                        value: appState.locale,
                        underline: SizedBox(),
                        items: [
                          DropdownMenuItem(value: 'en', child: Text('English')),
                          DropdownMenuItem(value: 'es', child: Text('Español')),
                        ],
                        onChanged: (value) {
                          if (value != null) {
                            appState.setLocale(value);
                            // Trigger rebuild
                            setState(() {});
                          }
                        },
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(height: 16),

              // Display Settings
              Card(
                child: Column(
                  children: [
                    SwitchListTile(
                      secondary: Icon(Icons.dark_mode),
                      title: Text(context.translate('settings_privacy')),
                      value: appState.isDarkMode,
                      onChanged: (value) {
                        appState.toggleDarkMode();
                      },
                    ),
                  ],
                ),
              ),
              SizedBox(height: 16),

              // Notifications
              Card(
                child: Column(
                  children: [
                    SwitchListTile(
                      secondary: Icon(Icons.notifications),
                      title: Text(context.translate('settings_notifications')),
                      value: appState.notificationsEnabled,
                      onChanged: (value) {
                        appState.toggleNotifications();
                      },
                    ),
                  ],
                ),
              ),
              SizedBox(height: 16),

              // About
              Card(
                child: ListTile(
                  leading: Icon(Icons.info),
                  title: Text(context.translate('settings_about')),
                  trailing: Icon(Icons.arrow_forward),
                  onTap: _showAboutDialog,
                ),
              ),
              SizedBox(height: 16),

              // Logout
              Card(
                child: ListTile(
                  leading: Icon(Icons.logout, color: AppTheme.errorColor),
                  title: Text(
                    context.translate('settings_logout'),
                    style: TextStyle(color: AppTheme.errorColor),
                  ),
                  onTap: _confirmLogout,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showAboutDialog() {
    showAboutDialog(
      context: context,
      applicationName: context.translate('app_name'),
      applicationVersion: '1.0.0',
      applicationLegalese: 'Copyright © 2024',
      children: [
        Padding(
          padding: EdgeInsets.only(top: 16),
          child: Text(
            context.translate('app_subtitle'),
          ),
        ),
      ],
    );
  }

  void _confirmLogout() {
    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(context.translate('settings_logout')),
          content: Text(context.translate('settings_confirm_logout')),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text(context.translate('common_cancel')),
            ),
            ElevatedButton(
              onPressed: () {
                Provider.of<AuthProvider>(context, listen: false).logout();
                Navigator.pushAndRemoveUntil(
                  context,
                  MaterialPageRoute(builder: (_) => LoginScreen()),
                  (route) => false,
                );
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: AppTheme.errorColor,
              ),
              child: Text(context.translate('common_yes')),
            ),
          ],
        );
      },
    );
  }
}
