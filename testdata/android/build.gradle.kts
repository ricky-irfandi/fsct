def localProperties = new File(localPropertiesFile).exists() ? new Properties() : new Properties()
def localPropertiesFile = rootProject.file('local.properties')
localProperties.load(new FileInputStream(localPropertiesFile))

def flutterVersionCode = localProperties['flutter.versionCode']
if (flutterVersionCode == null) {
    flutterVersionCode = '1'
}

def flutterVersionName = localProperties['flutter.versionName']
if (flutterVersionName == null) {
    flutterVersionName = '1.0.0'
}

plugins {
    id "com.android.application"
    id "kotlin-android"
    id "dev.flutter.flutter-gradle-plugin"
}

android {
    namespace "com.example.myapp"
    compileSdk 34
    ndkVersion "25.2.9519653"

    compileOptions {
        sourceCompatibility JavaVersion.VERSION_1_8
        targetCompatibility JavaVersion.VERSION_1_8
    }

    kotlinOptions {
        jvmTarget = "1.8"
    }

    sourceSets {
        main.java.srcDirs += 'src/main/kotlin'
    }

    defaultConfig {
        applicationId = "com.example.myapp"
        minSdkVersion = 21
        targetSdkVersion = 34
        versionCode = 1
        versionName = "1.0.0"
        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"
    }

    buildTypes {
        release {
            signingConfig = signingConfigs.debug
        }
    }

    lint {
        abortOnError = false
        checkReleaseBuilds = true
        warningsAsErrors = false
    }
}

dependencies {
    implementation "org.jetbrains.kotlin:kotlin-stdlib-jdk7:$kotlin_version"
    testImplementation "junit:junit:4.12"
    androidTestImplementation "androidx.test:runner:1.1.1"
    androidTestImplementation "androidx.test.espresso:espresso-core:3.1.1"
}
