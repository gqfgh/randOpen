<script setup>
import {CopyFileName, DelFile, OpenFile, RenameFile, SelectDir, TraversalFiles} from "../wailsjs/go/main/App.js";
import {ref} from "vue";
import {CaretRight, CopyDocument} from "@element-plus/icons-vue"

let dir = ref("h:/rather/vid")
const traversalFiles = () => TraversalFiles(dir.value)
const selectDir = () => SelectDir().then(result => dir.value = result)
</script>

<template>
  <el-space :direction="'vertical'" :size="'large'">
    <el-radio-group v-model="dir">
      <el-radio label="h:/rather/vid" size="large">视频</el-radio>
      <el-radio label="h:/rather/pic" size="large">图片</el-radio>
    </el-radio-group>
    <el-space>
      <el-button type="primary" @click="selectDir">选择文件夹</el-button>
      <el-input v-model="dir" disabled/>
    </el-space>
    <el-button @click="traversalFiles" :icon="CaretRight">遍历文件</el-button>
    <el-button-group>
      <el-button @click="OpenFile" color="#626aef">打开</el-button>
      <el-button @click="DelFile" color="#626aef">删除</el-button>
      <el-button @click="RenameFile" color="#626aef">重命名</el-button>
      <el-button @click="CopyFileName" :icon="CopyDocument" color="#626aef"/><!-- 复制文件名 -->
    </el-button-group>
  </el-space>
</template>