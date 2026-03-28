#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

// 读取 package.json 获取版本号
const packageJsonPath = path.join(__dirname, 'package.json');
const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
const version = packageJson.version;

console.log(`Syncing version: ${version}`);

// 更新 build.gradle
const buildGradlePath = path.join(__dirname, 'android', 'app', 'build.gradle');
let buildGradle = fs.readFileSync(buildGradlePath, 'utf8');

// 使用正则表达式替换 versionName
buildGradle = buildGradle.replace(
  /versionName ".*"/,
  `versionName "${version}"`
);

// 写回文件
fs.writeFileSync(buildGradlePath, buildGradle, 'utf8');
console.log(`Updated build.gradle to version ${version}`);

// 更新 version.json
const versionJsonPath = path.join(__dirname, 'public', 'version.json');
fs.writeFileSync(versionJsonPath, JSON.stringify({ version }, null, 2), 'utf-8');
console.log(`Updated version.json to version ${version}`);
