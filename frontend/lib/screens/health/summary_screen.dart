import 'package:flutter/material.dart';

import '../../config/localization.dart';
import '../../config/theme.dart';

class SummaryScreen extends StatefulWidget {
  @override
  State<SummaryScreen> createState() => _SummaryScreenState();
}

class _SummaryScreenState extends State<SummaryScreen> {
  int _selectedDays = 7;
  bool _isLoading = false;
  String? _summary;
  List<String>? _findings;
  String? _recommendations;

  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;

    return Scaffold(
      appBar: AppBar(
        title: Text(context.translate('health_summary_title')),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: EdgeInsets.all(isMobile ? 16 : 24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                context.translate('health_summary_title'),
                style: Theme.of(context).textTheme.headlineSmall,
              ),
              SizedBox(height: 24),

              // Days selector
              Text(
                'Period',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              SizedBox(height: 12),
              Wrap(
                spacing: 8,
                children: [7, 14, 30, 90].map((days) {
                  return FilterChip(
                    selected: _selectedDays == days,
                    label: Text('Last $days days'),
                    onSelected: (selected) {
                      setState(() => _selectedDays = days);
                    },
                  );
                }).toList(),
              ),
              SizedBox(height: 24),

              // Generate button
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _isLoading ? null : _generateSummary,
                  child: _isLoading
                      ? SizedBox(
                    height: 20,
                    width: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                      valueColor: AlwaysStoppedAnimation(Colors.white),
                    ),
                  )
                      : Text(context.translate('health_summary_generate')),
                ),
              ),
              SizedBox(height: 32),

              // Results
              if (_summary != null) ...[
                Card(
                  child: Padding(
                    padding: EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          context.translate('health_summary_last_days', {
                            'days': _selectedDays.toString(),
                          }),
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        SizedBox(height: 12),
                        Text(
                          _summary ?? '',
                          style: Theme.of(context).textTheme.bodyMedium,
                        ),
                      ],
                    ),
                  ),
                ),
                SizedBox(height: 16),
              ],

              if (_findings != null && _findings!.isNotEmpty) ...[
                Text(
                  context.translate('health_summary_key_findings'),
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                SizedBox(height: 12),
                ..._findings!.map((finding) => Card(
                  child: Padding(
                    padding: EdgeInsets.all(12),
                    child: Row(
                      children: [
                        Icon(Icons.check, color: AppTheme.successColor),
                        SizedBox(width: 12),
                        Expanded(child: Text(finding)),
                      ],
                    ),
                  ),
                )),
                SizedBox(height: 16),
              ],

              if (_recommendations != null) ...[
                Text(
                  context.translate('health_summary_recommendations'),
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                SizedBox(height: 12),
                Card(
                  child: Padding(
                    padding: EdgeInsets.all(16),
                    child: Text(
                      _recommendations ?? '',
                      style: Theme.of(context).textTheme.bodyMedium,
                    ),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Future<void> _generateSummary() async {
    setState(() => _isLoading = true);

    try {
      // TODO: Call gRPC to generate summary
      await Future.delayed(Duration(seconds: 2));

      setState(() {
        _summary = 'Health Summary for last $_selectedDays days: Overall health status is good with no critical findings.';
        _findings = [
          'Overall health status: Good',
          'Recent medications: None critical',
          'Blood pressure: Normal',
        ];
        _recommendations = 'Stay hydrated, maintain regular exercise, and schedule a check-up next month.';
        _isLoading = false;
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to generate summary')),
      );
      setState(() => _isLoading = false);
    }
  }
}
