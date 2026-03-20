<script setup lang="ts">
import { SaveOutline } from '@vicons/ionicons5'
import {
  NButton,
  NCard,
  NDatePicker,
  NDivider,
  NForm,
  NFormItem,
  NInput,
  NSkeleton,
  NSpace,
  NSwitch,
  NThing,
} from 'naive-ui'

import { ScrollContainer } from '@/components'
import MultiImageInput from '@/components/image-picker/MultiImageInput.vue'

import { useGalleryForm } from './composables/use-gallery-form'

defineOptions({ name: 'GalleryEdit' })

const { form, loading, saving, imageProcessing, isCreating, save } = useGalleryForm()
</script>

<template>
  <ScrollContainer wrapper-class="p-4 md:p-6 pb-10">
    <NCard :bordered="false" class="shadow-sm">
      <template #header>
        <NThing
          :title="isCreating ? '新建日常' : '编辑日常'"
          description="轻量记录文字与图片，不使用二级详情页。"
        />
      </template>

      <div v-if="loading" class="space-y-4 py-2">
        <NSkeleton text :repeat="6" />
      </div>

      <NForm v-else label-placement="top" class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_320px]">
        <NCard size="small" embedded>
          <NFormItem label="文案内容" path="content">
            <NInput
              v-model:value="form.content"
              type="textarea"
              :autosize="{ minRows: 10, maxRows: 18 }"
              placeholder="写下今天的光线、风声、城市、心情。支持换行，但不需要复杂 markdown。"
            />
          </NFormItem>

          <NDivider />

          <NFormItem label="图片列表">
            <MultiImageInput v-model:value="form.images" />
          </NFormItem>
        </NCard>

        <NCard size="small" embedded>
          <NFormItem label="发布时间">
            <NDatePicker
              v-model:value="form.createdAt"
              type="datetime"
              clearable
              class="w-full"
            />
          </NFormItem>

          <NFormItem label="发布状态">
            <NSwitch v-model:value="form.isPublished">
              <template #checked>已发布</template>
              <template #unchecked>草稿</template>
            </NSwitch>
          </NFormItem>

          <NFormItem label="置顶">
            <NSwitch v-model:value="form.isTop">
              <template #checked>置顶</template>
              <template #unchecked>普通</template>
            </NSwitch>
          </NFormItem>

          <div class="rounded-xl border border-current/10 bg-black/2 p-4 text-sm text-current/70 dark:bg-white/2">
            图片元信息处理中：{{ imageProcessing ? '进行中' : '空闲' }}
          </div>

          <template #footer>
            <NSpace justify="end">
              <NButton type="primary" :loading="saving" @click="save">
                <template #icon>
                  <SaveOutline />
                </template>
                {{ isCreating ? '创建日常' : '保存更新' }}
              </NButton>
            </NSpace>
          </template>
        </NCard>
      </NForm>
    </NCard>
  </ScrollContainer>
</template>