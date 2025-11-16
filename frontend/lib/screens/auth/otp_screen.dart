import 'package:flutter/material.dart';
import 'package:pin_code_fields/pin_code_fields.dart';
import 'package:provider/provider.dart';

import '../../config/localization.dart';
import '../../config/theme.dart';
import '../../providers/auth_provider.dart';
import '../home_screen.dart';

class OTPScreen extends StatefulWidget {
  final String email;

  const OTPScreen({required this.email});

  @override
  State<OTPScreen> createState() => _OTPScreenState();
}

class _OTPScreenState extends State<OTPScreen> {
  late TextEditingController _otpController;
  bool _isLoading = false;
  int _resendCountdown = 0;

  @override
  void initState() {
    super.initState();
    _otpController = TextEditingController();
    _startResendCountdown();
  }

  @override
  void dispose() {
    _otpController.dispose();
    super.dispose();
  }

  void _startResendCountdown() {
    _resendCountdown = 60;
    Future.doWhile(() async {
      await Future.delayed(Duration(seconds: 1));
      if (mounted) {
        setState(() => _resendCountdown--);
      }
      return _resendCountdown > 0;
    });
  }

  void _verifyOTP() async {
    if (_otpController.text.length != 6) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(context.translate('validation_otp_invalid')),
          backgroundColor: AppTheme.errorColor,
        ),
      );
      return;
    }

    setState(() => _isLoading = true);

    final authProvider = Provider.of<AuthProvider>(context, listen: false);
    final success =
        await authProvider.verifyOTP(widget.email, _otpController.text);

    if (mounted) {
      setState(() => _isLoading = false);

      if (success) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (_) => HomeScreen()),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content:
                Text(authProvider.error ?? context.translate('auth_error')),
            backgroundColor: AppTheme.errorColor,
          ),
        );
      }
    }
  }

  void _resendOTP() async {
    final authProvider = Provider.of<AuthProvider>(context, listen: false);
    final success = await authProvider.sendOTP(widget.email);

    if (success) {
      _otpController.clear();
      _startResendCountdown();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(context.translate('auth_otp_sent')),
          backgroundColor: AppTheme.successColor,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;

    return Scaffold(
      appBar: AppBar(
        title: Text(context.translate('auth_enter_otp')),
        leading: IconButton(
          icon: Icon(Icons.arrow_back),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: EdgeInsets.all(isMobile ? 24 : 48),
          child: Center(
            child: ConstrainedBox(
              constraints: BoxConstraints(maxWidth: 400),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.stretch,
                children: [
                  Icon(
                    Icons.verified_user,
                    size: 80,
                    color: AppTheme.primaryColor,
                  ),
                  SizedBox(height: 24),
                  Text(
                    context.translate('auth_enter_otp'),
                    textAlign: TextAlign.center,
                    style: Theme.of(context).textTheme.headlineSmall,
                  ),
                  SizedBox(height: 8),
                  Text(
                    'Sent to ${widget.email}',
                    textAlign: TextAlign.center,
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          color: AppTheme.textSecondary,
                        ),
                  ),
                  SizedBox(height: 48),
                  PinCodeTextField(
                    appContext: context,
                    length: 6,
                    controller: _otpController,
                    keyboardType: TextInputType.number,
                    pinTheme: PinTheme(
                      shape: PinCodeFieldShape.box,
                      borderRadius: BorderRadius.circular(8),
                      fieldHeight: 50,
                      fieldWidth: 50,
                      activeFillColor: AppTheme.surfaceColor,
                      inactiveFillColor: AppTheme.surfaceColor,
                      selectedFillColor: Colors.white,
                      activeColor: AppTheme.primaryColor,
                      inactiveColor: AppTheme.textHint,
                      selectedColor: AppTheme.primaryColor,
                    ),
                    enabled: !_isLoading,
                    onChanged: (_) {},
                  ),
                  SizedBox(height: 32),
                  ElevatedButton(
                    onPressed: _isLoading ? null : _verifyOTP,
                    child: _isLoading
                        ? SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor: AlwaysStoppedAnimation(Colors.white),
                            ),
                          )
                        : Text(context.translate('auth_verify_otp')),
                  ),
                  SizedBox(height: 24),
                  Center(
                    child: _resendCountdown > 0
                        ? Text(
                            context.translate('auth_resend_in', {
                              'seconds': _resendCountdown.toString(),
                            }),
                            style:
                                Theme.of(context).textTheme.bodySmall?.copyWith(
                                      color: AppTheme.textSecondary,
                                    ),
                          )
                        : TextButton(
                            onPressed: _resendOTP,
                            child: Text(context.translate('auth_resend_otp')),
                          ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
