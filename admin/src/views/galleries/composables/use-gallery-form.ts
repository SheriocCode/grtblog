import { useMessage } from 'naive-ui'
import { reactive, ref, computed, onMounted, toRef } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useLeaveConfirm } from '@/composables'
import { createGallery, getGallery, updateGallery } from '@/services/galleries'
import type { ContentExtInfo } from '@/types/ext-info'
import { useImageExtInfo } from '@/composables/use-image-ext-info'

function splitImages(value: string) {
  return value
    .split(/\r?\n/)
    .map((item) => item.trim())
    .filter(Boolean)
}

function joinImages(images?: string[]) {
  return (images ?? []).join('\n')
}

export function useGalleryForm() {
  const route = useRoute()
  const router = useRouter()
  const message = useMessage()

  const galleryId = computed(() => {
    const param = route.params.id
    if (!param || param === 'new') return null
    const id = Number(param)
    return Number.isFinite(id) ? id : null
  })

  const isCreating = computed(() => galleryId.value === null)
  const loading = ref(false)
  const saving = ref(false)
  const initialSnapshot = ref('')

  const form = reactive({
    content: '',
    images: '',
    isPublished: false,
    isTop: false,
    createdAt: Date.now() as number | null,
  })

  const baseExtInfo = ref<ContentExtInfo | null>(null)
  const { extInfo, processing } = useImageExtInfo({
    content: toRef(form, 'content'),
    extraImages: toRef(form, 'images'),
    baseExtInfo,
  })

  const takeSnapshot = () => JSON.stringify(form)
  const isDirty = computed(() => initialSnapshot.value !== '' && takeSnapshot() !== initialSnapshot.value)

  async function fetch() {
    if (isCreating.value) {
      initialSnapshot.value = takeSnapshot()
      return null
    }

    loading.value = true
    try {
      const data = await getGallery(galleryId.value!)
      form.content = data.content
      form.images = joinImages(data.images)
      form.isPublished = data.isPublished
      form.isTop = data.isTop
      form.createdAt = data.createdAt ? new Date(data.createdAt).getTime() : null
      baseExtInfo.value = data.extInfo ?? null
      initialSnapshot.value = takeSnapshot()
      return data
    } catch (error) {
      console.error(error)
      message.error('无法加载日常数据')
      router.replace({ name: 'galleryList' })
      return null
    } finally {
      loading.value = false
    }
  }

  async function save() {
    if (!form.content.trim()) return message.error('请输入文案内容')

    saving.value = true
    try {
      const payload = {
        content: form.content,
        images: splitImages(form.images),
        isPublished: form.isPublished,
        isTop: form.isTop,
        createdAt: form.createdAt ? new Date(form.createdAt).toISOString() : undefined,
        extInfo: extInfo.value ?? undefined,
      }

      if (isCreating.value) {
        await createGallery(payload)
        message.success('创建成功')
      } else {
        await updateGallery(galleryId.value!, payload)
        message.success('更新成功')
      }

      initialSnapshot.value = takeSnapshot()
      router.push({ name: 'galleryList' })
    } catch (error: any) {
      message.error(error.message || '保存失败')
    } finally {
      saving.value = false
    }
  }

  useLeaveConfirm({
    when: isDirty,
    title: '未保存的更改',
    content: '当前内容未保存，确定要离开吗？',
    positiveText: '离开',
    negativeText: '继续编辑',
  })

  onMounted(fetch)

  return {
    form,
    loading,
    saving,
    imageProcessing: processing,
    isCreating,
    isDirty,
    fetch,
    save,
  }
}