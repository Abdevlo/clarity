import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:uuid/uuid.dart';

import '../../config/localization.dart';
import '../../config/theme.dart';
import '../../providers/auth_provider.dart';

class AIHealthAssistantScreen extends StatefulWidget {
  @override
  State<AIHealthAssistantScreen> createState() => _AIHealthAssistantScreenState();
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

class _AIHealthAssistantScreenState extends State<AIHealthAssistantScreen> {
  late TextEditingController _messageController;
  late ScrollController _scrollController;
  final List<ChatMessage> _messages = [];
  bool _isLoading = false;
  final String _conversationId = const Uuid().v4();

  @override
  void initState() {
    super.initState();
    _messageController = TextEditingController();
    _scrollController = ScrollController();

    // Add welcome message
    _messages.add(ChatMessage(
      id: const Uuid().v4(),
      message: 'Hello! I\'m your AI Doctor. How can I help you today?',
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
      // TODO: Call gRPC to send message
      await Future.delayed(Duration(seconds: 2));

      final response = 'Thank you for sharing. Based on your symptoms, I recommend...';

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
        SnackBar(content: Text(context.translate('doctor_chat_error'))),
      );
      setState(() => _isLoading = false);
    }
  }

  void _scrollToBottom() {
    Future.delayed(Duration(milliseconds: 300), () {
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: Duration(milliseconds: 300),
        curve: Curves.easeOut,
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;

    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(context.translate('doctor_chat_title')),
            Text(
              context.translate('doctor_chat_ai_powered'),
              style: TextStyle(fontSize: 12, color: Colors.white70),
            ),
          ],
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: ListView.builder(
              controller: _scrollController,
              padding: EdgeInsets.symmetric(
                horizontal: isMobile ? 12 : 24,
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
                      hintText: context.translate('doctor_chat_message_hint'),
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
                : MediaQuery.of(context).size.width * 0.6,
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
