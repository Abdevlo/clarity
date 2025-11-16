import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../config/localization.dart';
import '../config/theme.dart';
import '../providers/auth_provider.dart';
import '../providers/health_provider.dart';
import 'health/records_screen.dart';
import 'health/scanner_screen.dart';
import 'health/summary_screen.dart';
import 'ai_health_assistant_screen.dart';
import 'doctor/doctor_mode_chat_screen.dart';
import 'settings_screen.dart';

class HomeScreen extends StatefulWidget {
  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _selectedIndex = 0;

  late List<Widget> _screens;

  @override
  void initState() {
    super.initState();
    _screens = [
      DashboardScreen(),
      RecordsScreen(),
      ScannerScreen(),
      AIHealthAssistantScreen(),
      SettingsScreen(),
    ];

    // Fetch health records
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final authProvider = Provider.of<AuthProvider>(context, listen: false);
      if (authProvider.userId != null) {
        Provider.of<HealthProvider>(context, listen: false)
            .fetchRecords(authProvider.userId!);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthProvider>(context);

    return Scaffold(
      body: _screens[_selectedIndex],
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          authProvider.setDoctorMode(true);
          Navigator.push(
            context,
            MaterialPageRoute(builder: (_) => const DoctorModeChatScreen()),
          );
        },
        tooltip: 'Enter Doctor Mode',
        child: const Icon(Icons.medical_services),
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _selectedIndex,
        type: BottomNavigationBarType.fixed,
        onTap: (index) {
          setState(() => _selectedIndex = index);
        },
        items: [
          BottomNavigationBarItem(
            icon: const Icon(Icons.home),
            label: context.translate('dashboard_health_overview'),
          ),
          BottomNavigationBarItem(
            icon: const Icon(Icons.article),
            label: context.translate('records_title'),
          ),
          BottomNavigationBarItem(
            icon: const Icon(Icons.camera_alt),
            label: context.translate('scanner_title'),
          ),
          BottomNavigationBarItem(
            icon: const Icon(Icons.chat),
            label: context.translate('doctor_chat_title'),
          ),
          BottomNavigationBarItem(
            icon: const Icon(Icons.settings),
            label: context.translate('settings_title'),
          ),
        ],
      ),
    );
  }
}

class DashboardScreen extends StatefulWidget {
  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;
    final authProvider = Provider.of<AuthProvider>(context);
    final healthProvider = Provider.of<HealthProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(context.translate('app_name')),
        elevation: 0,
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: EdgeInsets.all(isMobile ? 16 : 24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Greeting
              Text(
                context.translate('dashboard_greeting', {
                  'name': authProvider.userEmail ?? 'User',
                }),
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              SizedBox(height: 24),

              // Quick Actions
              Text(
                context.translate('dashboard_quick_actions'),
                style: Theme.of(context).textTheme.titleLarge,
              ),
              SizedBox(height: 12),
              _buildQuickActions(context, isMobile),
              SizedBox(height: 32),

              // Recent Records
              Text(
                context.translate('dashboard_recent_records'),
                style: Theme.of(context).textTheme.titleLarge,
              ),
              SizedBox(height: 12),
              healthProvider.isLoading
                  ? Center(child: CircularProgressIndicator())
                  : healthProvider.records.isEmpty
                  ? Text(context.translate('records_empty'))
                  : ListView.builder(
                shrinkWrap: true,
                physics: NeverScrollableScrollPhysics(),
                itemCount: healthProvider.records.length,
                itemBuilder: (context, index) {
                  final record = healthProvider.records[index];
                  return Card(
                    margin: EdgeInsets.only(bottom: 12),
                    child: ListTile(
                      title: Text(record.title),
                      subtitle: Text(record.recordType),
                      trailing: Icon(Icons.arrow_forward),
                    ),
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildQuickActions(BuildContext context, bool isMobile) {
    final actions = [
      ('records_add_new', Icons.add, () {
        // Navigate to add record
      }),
      ('scanner_title', Icons.camera_alt, () {
        // Navigate to scanner
      }),
      ('health_summary_title', Icons.assessment, () {
        // Navigate to summary
      }),
      ('doctor_chat_title', Icons.chat, () {
        // Navigate to chat
      }),
    ];

    return GridView.builder(
      shrinkWrap: true,
      physics: NeverScrollableScrollPhysics(),
      gridDelegate: SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: isMobile ? 2 : 4,
        childAspectRatio: 1,
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
      ),
      itemCount: actions.length,
      itemBuilder: (context, index) {
        final (label, icon, onTap) = actions[index];
        return InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(12),
          child: Card(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(icon, size: 32, color: AppTheme.primaryColor),
                SizedBox(height: 8),
                Text(
                  context.translate(label),
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.bodySmall,
                ),
              ],
            ),
          ),
        );
      },
    );
  }
}
