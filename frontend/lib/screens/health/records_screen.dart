import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../config/localization.dart';
import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../providers/health_provider.dart';

class RecordsScreen extends StatefulWidget {
  @override
  State<RecordsScreen> createState() => _RecordsScreenState();
}

class _RecordsScreenState extends State<RecordsScreen> {
  String _selectedType = 'all';

  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;
    final authProvider = Provider.of<AuthProvider>(context);
    final healthProvider = Provider.of<HealthProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(context.translate('records_title')),
        elevation: 0,
      ),
      body: SafeArea(
        child: Column(
          children: [
            // Filter tabs
            SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              padding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              child: Row(
                children: [
                  _buildFilterChip('all', 'All'),
                  _buildFilterChip('prescription', context.translate('records_prescription')),
                  _buildFilterChip('appointment', context.translate('records_appointment')),
                  _buildFilterChip('lab_result', context.translate('records_lab_result')),
                  _buildFilterChip('symptom', context.translate('records_symptom')),
                ],
              ),
            ),
            Divider(height: 1),
            // Records list
            Expanded(
              child: healthProvider.isLoading
                  ? Center(child: CircularProgressIndicator())
                  : healthProvider.records.isEmpty
                  ? Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(Icons.article, size: 64, color: AppTheme.textHint),
                    SizedBox(height: 16),
                    Text(context.translate('records_empty')),
                  ],
                ),
              )
                  : ListView.builder(
                padding: EdgeInsets.all(12),
                itemCount: healthProvider.records.length,
                itemBuilder: (context, index) {
                  final record = healthProvider.records[index];
                  if (_selectedType != 'all' && record.recordType != _selectedType) {
                    return SizedBox.shrink();
                  }
                  return _buildRecordCard(context, record);
                },
              ),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _showAddRecordDialog,
        child: Icon(Icons.add),
      ),
    );
  }

  Widget _buildFilterChip(String type, String label) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 4),
      child: FilterChip(
        selected: _selectedType == type,
        label: Text(label),
        onSelected: (selected) {
          setState(() => _selectedType = type);
        },
      ),
    );
  }

  Widget _buildRecordCard(BuildContext context, HealthRecord record) {
    final iconMap = {
      'prescription': Icons.medication,
      'appointment': Icons.calendar_today,
      'lab_result': Icons.science,
      'symptom': Icons.sick,
    };

    return Card(
      margin: EdgeInsets.symmetric(vertical: 8),
      child: ListTile(
        leading: Icon(
          iconMap[record.recordType] ?? Icons.article,
          color: AppTheme.primaryColor,
        ),
        title: Text(record.title),
        subtitle: Text(record.description),
        trailing: Icon(Icons.chevron_right),
        onTap: () {
          // Show record details
        },
      ),
    );
  }

  void _showAddRecordDialog() {
    final recordTypeController = TextEditingController();
    final titleController = TextEditingController();
    final descriptionController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: Text(context.translate('records_add_new')),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                DropdownButtonFormField<String>(
                  value: 'prescription',
                  items: ['prescription', 'appointment', 'lab_result', 'symptom']
                      .map((e) => DropdownMenuItem(value: e, child: Text(e)))
                      .toList(),
                  onChanged: (value) {
                    recordTypeController.text = value ?? '';
                  },
                ),
                SizedBox(height: 12),
                TextField(
                  controller: titleController,
                  decoration: InputDecoration(
                    labelText: context.translate('records_title'),
                  ),
                ),
                SizedBox(height: 12),
                TextField(
                  controller: descriptionController,
                  decoration: InputDecoration(
                    labelText: context.translate('common_save'),
                  ),
                  maxLines: 3,
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text(context.translate('common_cancel')),
            ),
            ElevatedButton(
              onPressed: () {
                final authProvider = Provider.of<AuthProvider>(context, listen: false);
                final healthProvider = Provider.of<HealthProvider>(context, listen: false);

                healthProvider.createRecord(
                  userId: authProvider.userId ?? '',
                  recordType: recordTypeController.text.isNotEmpty
                      ? recordTypeController.text
                      : 'prescription',
                  title: titleController.text,
                  description: descriptionController.text,
                  metadata: {},
                );

                Navigator.pop(context);
              },
              child: Text(context.translate('common_save')),
            ),
          ],
        );
      },
    );
  }
}
