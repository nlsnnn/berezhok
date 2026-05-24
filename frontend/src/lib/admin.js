export const ADMIN_ROLE_LABELS = {
  super_admin: 'Суперадминистратор',
  admin: 'Администратор',
  support: 'Поддержка',
}

export function buildAdminUser(rawUser = {}) {
  return {
    id: rawUser?.id || rawUser?.user_id || null,
    email: rawUser?.email || null,
    name: rawUser?.name || 'Администратор',
    role: rawUser?.role || 'support',
  }
}

export function canManageAdmins(user) {
  return user?.role === 'super_admin'
}

export function canMutateOperations(user) {
  return user?.role === 'super_admin' || user?.role === 'admin'
}

export function getAdminRoleLabel(role) {
  return ADMIN_ROLE_LABELS[role] || role || '—'
}
