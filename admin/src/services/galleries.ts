import { request } from './http'
import type { ContentExtInfo } from '@/types/ext-info'

export interface GalleryListItem {
  id: number
  content: string
  contentHash: string
  images?: string[]
  imageCount: number
  isPublished: boolean
  isTop: boolean
  createdAt: string
  updatedAt: string
}

export interface GalleryDetail extends GalleryListItem {
  authorId: number
  extInfo?: ContentExtInfo | null
}

export interface GalleryListResponse {
  items: GalleryListItem[]
  total: number
  page: number
  size: number
}

export interface ListGalleriesParams {
  page?: number
  pageSize?: number
  published?: boolean
}

export interface CreateGalleryPayload {
  content: string
  images?: string[]
  isPublished: boolean
  isTop: boolean
  extInfo?: ContentExtInfo | null
  createdAt?: string | null
}

export interface UpdateGalleryPayload extends CreateGalleryPayload {}

function stripEmpty<T extends object>(value: T): Record<string, unknown> {
  return Object.fromEntries(
    Object.entries(value).filter(
      ([, entry]) => entry !== undefined && entry !== null && entry !== '',
    ),
  )
}

export function listGalleries(params: ListGalleriesParams) {
  return request<GalleryListResponse>('/admin/galleries', {
    method: 'GET',
    query: stripEmpty(params),
  })
}

export function getGallery(id: number) {
  return request<GalleryDetail>(`/admin/galleries/${id}`, {
    method: 'GET',
  })
}

export function createGallery(payload: CreateGalleryPayload) {
  return request<GalleryDetail>('/galleries', {
    method: 'POST',
    body: payload,
  })
}

export function updateGallery(id: number, payload: UpdateGalleryPayload) {
  return request<GalleryDetail>(`/galleries/${id}`, {
    method: 'PUT',
    body: payload,
  })
}

export function deleteGallery(id: number) {
  return request<void>(`/galleries/${id}`, {
    method: 'DELETE',
  })
}

export function batchSetGalleryPublished(payload: { ids: number[]; isPublished: boolean }) {
  return request<void>('/admin/galleries/published', {
    method: 'PUT',
    body: payload,
  })
}

export function batchSetGalleryTop(payload: { ids: number[]; isTop: boolean }) {
  return request<void>('/admin/galleries/top', {
    method: 'PUT',
    body: payload,
  })
}

export function batchDeleteGalleries(payload: { ids: number[] }) {
  return request<void>('/admin/galleries/batch-delete', {
    method: 'POST',
    body: payload,
  })
}