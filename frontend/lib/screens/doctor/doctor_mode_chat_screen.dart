import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:uuid/uuid.dart';

import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../../providers/health_provider.dart';

class DoctorModeChatScreen extends StatefulWidget {
  const DoctorModeChatScreen({Key? key}) : super(key: key);

  @override
  State<DoctorModeChatScreen> createState() => _DoctorModeChatScreenState();
}

class ChatMessage {
  final String id;
  final String message;
  final bool isUserMessage;
  final DateTime timestamp;

  ChatMessage({
    required this.id,
    required this.message,
    required this.isUserMessage,
    required this.timestamp,
  });
}

class _DoctorModeChatScreenState extends State<DoctorModeChatScreen> {
  late TextEditingController _messageController;
  late ScrollController _scrollController;
  final List<ChatMessage> _messages = [];
  bool _isLoading = false;
  bool _showPatientDataPanel = true;

  @override
  void initState() {
    super.initState();
    _messageController = TextEditingController();
    _scrollController = ScrollController();

    // Add welcome message for doctor
    _messages.add(ChatMessage(
      id: const Uuid().v4(),
      message: 'Hello Doctor! I\'m your AI diagnostic assistant. I have access to your patient\'s health records, test results, and medical history. How can I help you analyze this patient\'s condition?',
      isUserMessage: false,
      timestamp: DateTime.now(),
    ));
  }

  @override
  void dispose() {
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _sendMessage() async {
    if (_messageController.text.trim().isEmpty) return;

    final userMessage = _messageController.text.trim();
    _messageController.clear();

    setState(() {
      _messages.add(ChatMessage(
        id: const Uuid().v4(),
        message: userMessage,
        isUserMessage: true,
        timestamp: DateTime.now(),
      ));
      _isLoading = true;
    });

    _scrollToBottom();

    try {
      // TODO: Call gRPC to send message with patient context
      await Future.delayed(Duration(seconds: 2));

      final response = 'Based on the patient\'s records and current presentation, I recommend...';

      setState(() {
        _messages.add(ChatMessage(
          id: const Uuid().v4(),
          message: response,
          isUserMessage: false,
          timestamp: DateTime.now(),
        ));
        _isLoading = false;
      });

      _scrollToBottom();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error: ${e.toString()}')),
      );
      setState(() => _isLoading = false);
    }
  }

  void _scrollToBottom() {
    Future.delayed(Duration(milliseconds: 300), () {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;
    final healthProvider = Provider.of<HealthProvider>(context);
    final authProvider = Provider.of<AuthProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('AI Diagnostic Assistant'),
            Text(
              'Doctor Mode - Analyzing Patient Data',
              style: TextStyle(fontSize: 12, color: Colors.white70),
            ),
          ],
        ),
        elevation: 0,
        actions: [
          IconButton(
            icon: Icon(_showPatientDataPanel ? Icons.visibility : Icons.visibility_off),
            onPressed: () {
              setState(() => _showPatientDataPanel = !_showPatientDataPanel);
            },
            tooltip: _showPatientDataPanel ? 'Hide patient data' : 'Show patient data',
          ),
          IconButton(
            icon: const Icon(Icons.exit_to_app),
            onPressed: () {
              authProvider.setDoctorMode(false);
              Navigator.pop(context);
            },
            tooltip: 'Exit doctor mode',
          ),
        ],
      ),
      body: Row(
        children: [
          // Chat interface
          Expanded(
            flex: isMobile ? 1 : 2,
            child: Column(
              children: [
                Expanded(
                  child: ListView.builder(
                    controller: _scrollController,
                    padding: EdgeInsets.symmetric(
                      horizontal: isMobile ? 12 : 16,
                      vertical: 16,
                    ),
                    itemCount: _messages.length + (_isLoading ? 1 : 0),
                    itemBuilder: (context, index) {
                      if (index == _messages.length && _isLoading) {
                        return Padding(
                          padding: EdgeInsets.symmetric(vertical: 8),
                          child: Align(
                            alignment: Alignment.centerLeft,
                            child: Padding(
                              padding: EdgeInsets.only(left: 8),
                              child: SizedBox(
                                width: 30,
                                height: 30,
                                child: CircularProgressIndicator(strokeWidth: 2),
                              ),
                            ),
                          ),
                        );
                      }

                      final message = _messages[index];
                      return _buildMessageBubble(context, message, isMobile);
                    },
                  ),
                ),
                Divider(height: 1),
                Padding(
                  padding: EdgeInsets.all(12),
                  child: Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: _messageController,
                          enabled: !_isLoading,
                          decoration: InputDecoration(
                            hintText: 'Ask about patient diagnosis, treatment options...',
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(24),
                            ),
                            contentPadding: EdgeInsets.symmetric(
                              horizontal: 16,
                              vertical: 12,
                            ),
                          ),
                          maxLines: null,
                          textInputAction: TextInputAction.send,
                          onSubmitted: (_) => _sendMessage(),
                        ),
                      ),
                      SizedBox(width: 8),
                      FloatingActionButton(
                        mini: true,
                        onPressed: _isLoading ? null : _sendMessage,
                        child: Icon(Icons.send),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          // Patient data panel (only on desktop/tablet)
          if (!isMobile && _showPatientDataPanel)
            Container(
              width: 300,
              decoration: BoxDecoration(
                border: Border(left: BorderSide(color: Colors.grey[300]!)),
              ),
              child: _buildPatientDataPanel(context, healthProvider),
            ),
        ],
      ),
    );
  }

  Widget _buildPatientDataPanel(BuildContext context, HealthProvider healthProvider) {
    return SingleChildScrollView(
      padding: EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Patient Data',
            style: Theme.of(context).textTheme.titleLarge,
          ),
          SizedBox(height: 16),

          // Health Summary
          Card(
            child: Padding(
              padding: EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Health Summary',
                    style: Theme.of(context).textTheme.titleSmall,
                  ),
                  SizedBox(height: 8),
                  Text(
                    'Total Records: ${healthProvider.records.length}',
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                  Text(
                    'Last Updated: ${DateTime.now().toString().split('.')[0]}',
                    style: Theme.of(context).textTheme.bodySmall,
                  ),
                ],
              ),
            ),
          ),
          SizedBox(height: 16),

          // Recent Records
          Text(
            'Recent Records',
            style: Theme.of(context).textTheme.titleSmall,
          ),
          SizedBox(height: 8),
          if (healthProvider.records.isEmpty)
            Padding(
              padding: EdgeInsets.all(16),
              child: Text(
                'No records available',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: AppTheme.textSecondary,
                ),
              ),
            )
          else
            ListView.builder(
              shrinkWrap: true,
              physics: NeverScrollableScrollPhysics(),
              itemCount: healthProvider.records.take(5).length,
              itemBuilder: (context, index) {
                final record = healthProvider.records[index];
                return Card(
                  margin: EdgeInsets.only(bottom: 8),
                  child: Padding(
                    padding: EdgeInsets.all(8),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          record.title,
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                        SizedBox(height: 4),
                        Text(
                          record.recordType,
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: AppTheme.textSecondary,
                            fontSize: 11,
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
        ],
      ),
    );
  }

  Widget _buildMessageBubble(
    BuildContext context,
    ChatMessage message,
    bool isMobile,
  ) {
    return Padding(
      padding: EdgeInsets.symmetric(vertical: 8),
      child: Align(
        alignment: message.isUserMessage
            ? Alignment.centerRight
            : Alignment.centerLeft,
        child: ConstrainedBox(
          constraints: BoxConstraints(
            maxWidth: isMobile
                ? MediaQuery.of(context).size.width * 0.8
                : MediaQuery.of(context).size.width * 0.5,
          ),
          child: Container(
            padding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: BoxDecoration(
              color: message.isUserMessage
                  ? AppTheme.primaryColor
                  : AppTheme.surfaceColor,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Text(
              message.message,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: message.isUserMessage ? Colors.white : AppTheme.textPrimary,
              ),
            ),
          ),
        ),
      ),
    );
  }
}
