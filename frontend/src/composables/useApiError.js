import { useI18n } from '@/i18n/index.js'

export function useApiError() {
  const { t } = useI18n()

  /**
   * Resolves a structured { title, detail } from an Axios error.
   * Uses translated strings when the error_code is known,
   * falls back to the API's own title/details for unknown codes.
   */
  function resolveApiError(e) {
    const data = e?.response?.data
    if (!data?.error_code) {
      return { title: data?.title || t('error'), detail: data?.details || '' }
    }
    const code = data.error_code
    const titleKey = `err_${code}_title`
    const detailKey = `err_${code}_detail`
    const title  = t(titleKey)  !== titleKey  ? t(titleKey)  : (data.title   || t('error'))
    const detail = t(detailKey) !== detailKey ? t(detailKey) : (data.details || '')
    return { title, detail }
  }

  return { resolveApiError }
}