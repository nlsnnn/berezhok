export function toggleSelectedId(selectedIds, id) {
  if (selectedIds.includes(id)) {
    return selectedIds.filter((item) => item !== id)
  }
  return [...selectedIds, id]
}

export function toggleAllPageIds(selectedIds, pageIds) {
  const selectedSet = new Set(selectedIds)
  const allSelected = pageIds.length > 0 && pageIds.every((id) => selectedSet.has(id))

  if (allSelected) {
    return selectedIds.filter((id) => !pageIds.includes(id))
  }

  return Array.from(new Set([...selectedIds, ...pageIds]))
}

export function pruneSelectedIds(selectedIds, availableIds) {
  const availableSet = new Set(availableIds)
  return selectedIds.filter((id) => availableSet.has(id))
}
