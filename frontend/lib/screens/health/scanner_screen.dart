import 'dart:io';

import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';

import '../../config/localization.dart';
import '../../config/theme.dart';

class ScannerScreen extends StatefulWidget {
  @override
  State<ScannerScreen> createState() => _ScannerScreenState();
}

class _ScannerScreenState extends State<ScannerScreen> {
  XFile? _scannedImage;
  Map<String, String>? _extractedData;
  bool _isLoading = false;

  final ImagePicker _imagePicker = ImagePicker();

  @override
  Widget build(BuildContext context) {
    final isMobile = MediaQuery.of(context).size.width < 600;

    return Scaffold(
      appBar: AppBar(
        title: Text(context.translate('scanner_title')),
        elevation: 0,
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          child: Padding(
            padding: EdgeInsets.all(isMobile ? 16 : 24),
            child: _extractedData == null
                ? _buildScannerView(context, isMobile)
                : _buildResultsView(context, isMobile),
          ),
        ),
      ),
    );
  }

  Widget _buildScannerView(BuildContext context, bool isMobile) {
    return Column(
      children: [
        Container(
          height: 300,
          decoration: BoxDecoration(
            border: Border.all(color: AppTheme.textHint, width: 2),
            borderRadius: BorderRadius.circular(12),
            color: AppTheme.surfaceColor,
          ),
          child: _scannedImage == null
              ? Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.camera_alt, size: 48, color: AppTheme.textHint),
                SizedBox(height: 16),
                Text(
                  context.translate('scanner_upload'),
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                    color: AppTheme.textSecondary,
                  ),
                ),
              ],
            ),
          )
              : Image.file(
            File(_scannedImage!.path),
            fit: BoxFit.cover,
          ),
        ),
        SizedBox(height: 32),
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            onPressed: _isLoading ? null : _pickImage,
            child: _isLoading
                ? SizedBox(
              height: 20,
              width: 20,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                valueColor: AlwaysStoppedAnimation(Colors.white),
              ),
            )
                : Text(context.translate('scanner_upload')),
          ),
        ),
        SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton(
            onPressed: _isLoading ? null : _takePhoto,
            child: Text(context.translate('scanner_camera')),
          ),
        ),
      ],
    );
  }

  Widget _buildResultsView(BuildContext context, bool isMobile) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Card(
          child: Padding(
            padding: EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  context.translate('scanner_results'),
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                SizedBox(height: 16),
                ..._extractedData!.entries.map((e) => Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      e.key.replaceAll('_', ' ').toUpperCase(),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: AppTheme.textSecondary,
                      ),
                    ),
                    SizedBox(height: 4),
                    Text(
                      e.value,
                      style: Theme.of(context).textTheme.bodyLarge,
                    ),
                    SizedBox(height: 12),
                  ],
                )),
              ],
            ),
          ),
        ),
        SizedBox(height: 24),
        SizedBox(
          width: double.infinity,
          child: ElevatedButton(
            onPressed: () {
              // Save to records
            },
            child: Text(context.translate('scanner_save')),
          ),
        ),
        SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: OutlinedButton(
            onPressed: () {
              setState(() {
                _extractedData = null;
                _scannedImage = null;
              });
            },
            child: Text(context.translate('common_back')),
          ),
        ),
      ],
    );
  }

  Future<void> _pickImage() async {
    try {
      final image = await _imagePicker.pickImage(source: ImageSource.gallery);
      if (image != null) {
        setState(() => _scannedImage = image);
        await _scanPrescription();
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(context.translate('scanner_error'))),
      );
    }
  }

  Future<void> _takePhoto() async {
    try {
      final image = await _imagePicker.pickImage(source: ImageSource.camera);
      if (image != null) {
        setState(() => _scannedImage = image);
        await _scanPrescription();
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(context.translate('scanner_error'))),
      );
    }
  }

  Future<void> _scanPrescription() async {
    setState(() => _isLoading = true);

    try {
      // TODO: Call gRPC to scan prescription
      await Future.delayed(Duration(seconds: 2));

      // Mock extracted data
      setState(() {
        _extractedData = {
          'medication': 'Aspirin',
          'dosage': '500mg',
          'frequency': 'Twice daily',
          'duration': '7 days',
          'indication': 'Headache relief',
        };
        _isLoading = false;
      });
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text(context.translate('scanner_error'))),
      );
      setState(() => _isLoading = false);
    }
  }
}
