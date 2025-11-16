import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class AppLocalizations {
  late Map<String, String> _localizedStrings;
  late Locale _locale;

  AppLocalizations(this._locale);

  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  Future<bool> load() async {
    try {
      final jsonString = await rootBundle.loadString(
        'assets/translations/${_locale.languageCode}.json',
      );
      final jsonMap = json.decode(jsonString);
      _localizedStrings = jsonMap.cast<String, String>();
      return true;
    } catch (e) {
      print('Error loading translations: $e');
      // Fallback to English
      final jsonString = await rootBundle.loadString(
        'assets/translations/en.json',
      );
      final jsonMap = json.decode(jsonString);
      _localizedStrings = jsonMap.cast<String, String>();
      return true;
    }
  }

  String translate(String key, [Map<String, String>? replacements]) {
    String value = _localizedStrings[key] ?? key;

    if (replacements != null) {
      replacements.forEach((k, v) {
        value = value.replaceAll('{$k}', v);
      });
    }

    return value;
  }

  String operator [](String key) => translate(key);
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) {
    return ['en', 'es'].contains(locale.languageCode);
  }

  @override
  Future<AppLocalizations> load(Locale locale) async {
    AppLocalizations localizations = AppLocalizations(locale);
    await localizations.load();
    return localizations;
  }

  @override
  bool shouldReload(LocalizationsDelegate<AppLocalizations> old) => false;
}

// Extension for easier access
extension AppLocalizationsX on BuildContext {
  AppLocalizations get locale => AppLocalizations.of(this);

  String translate(String key, [Map<String, String>? replacements]) =>
      AppLocalizations.of(this).translate(key, replacements);
}
