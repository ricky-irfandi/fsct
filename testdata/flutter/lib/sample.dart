import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;
  final String apiKey = "sk-test123456789";
  final String password = "password123";

  void _incrementCounter() {
    print("Counter incremented: $_counter");
    setState(() {
      _counter++;
    });
  }

  Future<void> _fetchData() async {
    final response = await http.get(Uri.parse('http://api.example.com/data'));
    print('Response status: ${response.statusCode}');
  }

  void _openUrl() async {
    final url = Uri.parse('https://example.com');
    if (await canLaunchUrl(url)) {
      await launchUrl(url);
    }
  }

  void _showPrivacyPolicy() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Privacy Policy'),
        content: const Text('This is our privacy policy.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }

  void _handleLogin() {
    print("Login clicked");
  }

  void _handleLogout() {
    print("Logout clicked");
  }

  void _handleDeleteAccount() {
    print("Delete account clicked");
  }

  void _handleForgotPassword() {
    print("Forgot password clicked");
  }

  bool signOut() {
    return true;
  }

  bool deleteAccount() {
    return true;
  }

  bool forgotPassword() {
    return true;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const Text(
              'You have pushed the button this many times:',
            ),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headlineMedium,
            ),
            ElevatedButton(
              onPressed: _openUrl,
              child: const Text('Open URL'),
            ),
            ElevatedButton(
              onPressed: _showPrivacyPolicy,
              child: const Text('Privacy Policy'),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ),
    );
  }
}
