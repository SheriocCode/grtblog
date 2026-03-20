import {
  NButton,
  NCard,
  NDataTable,
  NDropdown,
  NPagination,
  NPopconfirm,
  NSpace,
  NTag,
} from 'naive-ui'
import { defineComponent, ref } from 'vue'
import { useRouter } from 'vue-router'

import { ScrollContainer } from '@/components'
import { useTable } from '@/composables/table/use-table'
import { useDiscreteApi } from '@/composables/useDiscreteApi'
import {
  batchDeleteGalleries,
  batchSetGalleryPublished,
  batchSetGalleryTop,
  deleteGallery,
  listGalleries,
} from '@/services/galleries'

import type { GalleryListItem } from '@/services/galleries'
import type { DataTableColumns, DataTableRowKey } from 'naive-ui'

function excerpt(content: string, max = 72) {
  const normalized = content.replace(/\s+/g, ' ').trim()
  if (normalized.length <= max) return normalized || '-'
  return `${normalized.slice(0, max)}...`
}

export default defineComponent({
  name: 'GalleryList',
  setup() {
    const router = useRouter()
    const { message } = useDiscreteApi()
    const { data, loading, pagination, refresh } = useTable<GalleryListItem>(listGalleries)
    const checkedRowKeys = ref<DataTableRowKey[]>([])

    const handleCreate = () => {
      router.push({ name: 'galleryCreate' })
    }

    const handleEdit = (id: number) => {
      router.push({ name: 'galleryEdit', params: { id } })
    }

    const handleDelete = async (id: number) => {
      try {
        await deleteGallery(id)
        message.success('删除成功')
        refresh()
      } catch (error) {
        message.error(error instanceof Error ? error.message : '删除失败')
      }
    }

    const handleCheck = (rowKeys: DataTableRowKey[]) => {
      checkedRowKeys.value = rowKeys
    }

    const handleTogglePublished = async (row: GalleryListItem) => {
      try {
        await batchSetGalleryPublished({ ids: [row.id], isPublished: !row.isPublished })
        row.isPublished = !row.isPublished
        message.success(row.isPublished ? '已发布' : '已设为草稿')
      } catch (error) {
        message.error(error instanceof Error ? error.message : '操作失败')
      }
    }

    const handleToggleTop = async (row: GalleryListItem) => {
      try {
        await batchSetGalleryTop({ ids: [row.id], isTop: !row.isTop })
        row.isTop = !row.isTop
        message.success(row.isTop ? '已置顶' : '已取消置顶')
      } catch (error) {
        message.error(error instanceof Error ? error.message : '操作失败')
      }
    }

    const handleBatchPublish = async (isPublished: boolean) => {
      const ids = checkedRowKeys.value as number[]
      if (ids.length === 0) return
      try {
        await batchSetGalleryPublished({ ids, isPublished })
        data.value.forEach((item) => {
          if (ids.includes(item.id)) item.isPublished = isPublished
        })
        checkedRowKeys.value = []
        message.success(isPublished ? '批量发布成功' : '批量取消发布成功')
      } catch (error) {
        message.error(error instanceof Error ? error.message : '操作失败')
      }
    }

    const handleBatchTop = async (isTop: boolean) => {
      const ids = checkedRowKeys.value as number[]
      if (ids.length === 0) return
      try {
        await batchSetGalleryTop({ ids, isTop })
        data.value.forEach((item) => {
          if (ids.includes(item.id)) item.isTop = isTop
        })
        checkedRowKeys.value = []
        message.success(isTop ? '批量置顶成功' : '批量取消置顶成功')
      } catch (error) {
        message.error(error instanceof Error ? error.message : '操作失败')
      }
    }

    const handleBatchDelete = async () => {
      const ids = checkedRowKeys.value as number[]
      if (ids.length === 0) return
      try {
        await batchDeleteGalleries({ ids })
        checkedRowKeys.value = []
        message.success('批量删除成功')
        refresh()
      } catch (error) {
        message.error(error instanceof Error ? error.message : '操作失败')
      }
    }

    const batchPublishOptions = [
      { label: '设为已发布', key: 'publish' },
      { label: '设为草稿', key: 'unpublish' },
    ]

    const batchTopOptions = [
      { label: '设为置顶', key: 'top' },
      { label: '取消置顶', key: 'untop' },
    ]

    const columns: DataTableColumns<GalleryListItem> = [
      { type: 'selection' },
      {
        title: '内容',
        key: 'content',
        minWidth: 360,
        render: (row) => (
          <div class='space-y-1'>
            <div class='font-medium text-gray-700 dark:text-gray-200'>{excerpt(row.content, 96)}</div>
            <div class='text-xs text-gray-400'>Hash: {row.contentHash}</div>
          </div>
        ),
      },
      {
        title: '图片',
        key: 'imageCount',
        width: 100,
        render: (row) => <span class='font-mono text-xs text-gray-500'>{row.imageCount}</span>,
      },
      {
        title: '发布',
        key: 'isPublished',
        width: 100,
        render: (row) => (
          <span style={{ cursor: 'pointer' }} onClick={() => handleTogglePublished(row)}>
            <NTag size='small' type={row.isPublished ? 'success' : 'default'} bordered={false}>
              {row.isPublished ? '已发布' : '草稿'}
            </NTag>
          </span>
        ),
      },
      {
        title: '置顶',
        key: 'isTop',
        width: 100,
        render: (row) => (
          <span style={{ cursor: 'pointer' }} onClick={() => handleToggleTop(row)}>
            <NTag size='small' type={row.isTop ? 'warning' : 'default'} bordered={false}>
              {row.isTop ? '置顶' : '普通'}
            </NTag>
          </span>
        ),
      },
      {
        title: '发布时间',
        key: 'createdAt',
        width: 180,
        render: (row) => new Date(row.createdAt).toLocaleString(),
      },
      {
        title: '更新时间',
        key: 'updatedAt',
        width: 180,
        render: (row) => new Date(row.updatedAt).toLocaleString(),
      },
      {
        title: '操作',
        key: 'actions',
        width: 160,
        fixed: 'right',
        render: (row) => (
          <NSpace>
            <NButton size='small' type='primary' secondary onClick={() => handleEdit(row.id)}>
              编辑
            </NButton>
            <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
              {{
                trigger: () => <NButton size='small' type='error' secondary>删除</NButton>,
                default: () => '确认删除这条日常吗？',
              }}
            </NPopconfirm>
          </NSpace>
        ),
      },
    ]

    return () => (
      <div class='p-4 md:p-6'>
        <NCard bordered={false} class='shadow-sm'>
          {{
            header: () => '日常列表',
            headerExtra: () => (
              <NSpace>
                <NDropdown options={batchPublishOptions} onSelect={(key) => handleBatchPublish(key === 'publish')}>
                  <NButton disabled={checkedRowKeys.value.length === 0}>批量发布</NButton>
                </NDropdown>
                <NDropdown options={batchTopOptions} onSelect={(key) => handleBatchTop(key === 'top')}>
                  <NButton disabled={checkedRowKeys.value.length === 0}>批量置顶</NButton>
                </NDropdown>
                <NPopconfirm onPositiveClick={handleBatchDelete}>
                  {{
                    trigger: () => (
                      <NButton disabled={checkedRowKeys.value.length === 0} type='error' secondary>
                        批量删除
                      </NButton>
                    ),
                    default: () => '确认删除选中的日常吗？',
                  }}
                </NPopconfirm>
                <NButton type='primary' onClick={handleCreate}>新建日常</NButton>
              </NSpace>
            ),
            default: () => (
              <ScrollContainer class='overflow-hidden'>
                <NDataTable
                  remote={false}
                  rowKey={(row: GalleryListItem) => row.id}
                  columns={columns}
                  data={data.value}
                  loading={loading.value}
                  checkedRowKeys={checkedRowKeys.value}
                  onUpdateCheckedRowKeys={handleCheck}
                  maxHeight={640}
                />
                <div class='mt-4 flex justify-end'>
                  <NPagination
                    page={pagination.page}
                    pageSize={pagination.pageSize}
                    itemCount={pagination.itemCount}
                    showSizePicker={pagination.showSizePicker}
                    pageSizes={pagination.pageSizes}
                    onUpdatePage={pagination.onChange}
                    onUpdatePageSize={pagination.onUpdatePageSize}
                  />
                </div>
              </ScrollContainer>
            ),
          }}
        </NCard>
      </div>
    )
  },
})