import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dart:ui';

import '../config/localization.dart';
import '../config/theme.dart';
import '../providers/app_state_provider.dart';
import 'auth/login_screen.dart';

class OnboardingScreen extends StatefulWidget {
  @override
  State<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends State<OnboardingScreen> {
  late PageController _pageController;
  int _currentPage = 0;

  final List<OnboardingPage> pages = [
    OnboardingPage(
      title: 'onboarding_welcome_title',
      description: 'onboarding_welcome_desc',
      icon: Icons.health_and_safety,
    ),
    OnboardingPage(
      title: 'onboarding_secure_title',
      description: 'onboarding_secure_desc',
      icon: Icons.lock,
    ),
    OnboardingPage(
      title: 'onboarding_ai_title',
      description: 'onboarding_ai_desc',
      icon: Icons.smart_toy,
    ),
    OnboardingPage(
      title: 'onboarding_doctor_title',
      description: 'onboarding_doctor_desc',
      icon: Icons.chat,
    ),
  ];

  @override
  void initState() {
    super.initState();
    _pageController = PageController();
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  void _completeOnboarding() {
    final appState = Provider.of<AppStateProvider>(context, listen: false);
    appState.completeOnboarding();

    Navigator.of(context).pushReplacement(
      MaterialPageRoute(builder: (_) => LoginScreen()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          image: DecorationImage(
            image: AssetImage('assets/eyes.webp'),
            fit: BoxFit.cover,
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              Expanded(
                child: PageView.builder(
                  controller: _pageController,
                  onPageChanged: (index) {
                    setState(() => _currentPage = index);
                  },
                  itemCount: pages.length,
                  itemBuilder: (context, index) {
                    return OnboardingPageWidget(page: pages[index]);
                  },
                ),
              ),
              Padding(
                padding: EdgeInsets.symmetric(horizontal: 24, vertical: 24),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(50),
                  child: BackdropFilter(
                    filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
                    child: Container(
                      decoration: BoxDecoration(
                        color: Colors.white.withOpacity(0.15),
                      ),
                      padding: EdgeInsets.symmetric(horizontal: 10),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          // Left side: Arrow buttons
                          Row(
                            children: [
                              // Previous button
                              Material(
                                color: Colors.transparent,
                                child: InkWell(
                                  onTap: _currentPage > 0
                                      ? () {
                                          _pageController.previousPage(
                                            duration:
                                                Duration(milliseconds: 300),
                                            curve: Curves.easeInOut,
                                          );
                                        }
                                      : null,
                                  borderRadius: BorderRadius.circular(20),
                                  child: Container(
                                    padding: EdgeInsets.all(8),
                                    child: Icon(
                                      Icons.arrow_back,
                                      color: _currentPage > 0
                                          ? Colors.black
                                          : Colors.grey[400],
                                      size: 24,
                                    ),
                                  ),
                                ),
                              ),
                              SizedBox(width: 8),
                              // Next button
                              Material(
                                color: Colors.transparent,
                                child: InkWell(
                                  onTap: () {
                                    if (_currentPage == pages.length - 1) {
                                      _completeOnboarding();
                                    } else {
                                      _pageController.nextPage(
                                        duration: Duration(milliseconds: 300),
                                        curve: Curves.easeInOut,
                                      );
                                    }
                                  },
                                  borderRadius: BorderRadius.circular(20),
                                  child: Container(
                                    padding: EdgeInsets.all(8),
                                    child: Icon(
                                      Icons.arrow_forward,
                                      color: Colors.black,
                                      size: 24,
                                    ),
                                  ),
                                ),
                              ),
                            ],
                          ),
                          // Right side: Page indicators
                          Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: List.generate(
                              pages.length,
                              (index) => Container(
                                height: 8,
                                width: _currentPage == index ? 24 : 8,
                                margin: EdgeInsets.symmetric(horizontal: 4),
                                decoration: BoxDecoration(
                                  borderRadius: BorderRadius.circular(4),
                                  color: _currentPage == index
                                      ? const Color.fromARGB(255, 0, 242, 255)
                                      : Colors.white70,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class OnboardingPage {
  final String title;
  final String description;
  final IconData icon;

  OnboardingPage({
    required this.title,
    required this.description,
    required this.icon,
  });
}

class OnboardingPageWidget extends StatelessWidget {
  final OnboardingPage page;

  const OnboardingPageWidget({required this.page});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            page.icon,
            size: 96,
            color: AppTheme.primaryColor,
          ),
          SizedBox(height: 32),
          Text(
            context.translate(page.title),
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.displaySmall,
          ),
          SizedBox(height: 16),
          Text(
            context.translate(page.description),
            textAlign: TextAlign.center,
            style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                  color: AppTheme.textSecondary,
                ),
          ),
        ],
      ),
    );
  }
}
