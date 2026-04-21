'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from "@/contexts/AuthContext";
import { CreateEmployeeRequest, Employee, apiClient } from '@/lib/api';
import { AppLoading } from '@/components/ui/AppLoading';

export default function EmployeesPage() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalEmployees, setTotalEmployees] = useState(0);
  const [searchTerm, setSearchTerm] = useState('');
  
  const employeesPerPage = 10;
  
  const router = useRouter();
  const { user, logout } = useAuth();

  useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
    fetchEmployees();
  }, [user, currentPage, router]);

  const fetchEmployees = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const offset = (currentPage - 1) * employeesPerPage;
      const { data, error: apiError } = await apiClient.getEmployees(employeesPerPage, offset);
      
      if (data && !apiError) {
        setEmployees(data.employees || []);
        setTotalEmployees(data.total || 0);
      } else {
        setError(apiError?.message || 'Não foi possível carregar os colaboradores');
      }
    } catch (err) {
      setError('Erro de rede. Tente novamente.');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Tem a certeza de que pretende eliminar este colaborador?')) {
      return;
    }

    try {
      const { error } = await apiClient.deleteEmployee(id);
      if (!error) {
        setEmployees(employees.filter(emp => emp.id !== id));
        setTotalEmployees(prev => prev - 1);
      } else {
        alert('Não foi possível eliminar o colaborador: ' + error.message);
      }
    } catch (err) {
      alert('Erro de rede. Tente novamente.');
    }
  };

  const filteredEmployees = employees.filter(employee =>
    employee.first_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    employee.last_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    employee.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
    employee.department.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const totalPages = Math.ceil(totalEmployees / employeesPerPage);

  if (loading && employees.length === 0) {
    return (
      <div className="loadingScreen" aria-busy="true">
        <AppLoading size="lg" aria-busy={false} label="A carregar colaboradores" />
      </div>
    );
  }

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: 'var(--surface-page)',
    }}>
      {/* Header */}
      <header style={{
        backgroundColor: 'var(--surface-header)',
        boxShadow: 'var(--shadow-sm)',
        borderBottom: '1px solid var(--color-gray-200)',
      }}>
        <div style={{
          maxWidth: '1200px',
          margin: '0 auto',
          padding: '0 var(--space-4)',
        }}>
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            height: '64px',
          }}>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: 'var(--space-3)',
            }}>
              <div style={{
                width: '32px',
                height: '32px',
                backgroundColor: 'var(--color-primary)',
                borderRadius: 'var(--radius)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}>
                <svg 
                  style={{ width: '16px', height: '16px', color: 'var(--text-on-primary)' }} 
                  fill="none" 
                  viewBox="0 0 24 24" 
                  stroke="currentColor"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-2m-2 0H7m5 0v-5a2 2 0 012-2h2a2 2 0 012 2v5" />
                </svg>
              </div>
              <div>
                <h1 style={{
                  fontSize: '1.25rem',
                  fontWeight: '700',
                  color: 'var(--color-gray-900)',
                  margin: 0,
                }}>
                  GonsGarage
                </h1>
                <p style={{
                  fontSize: '0.75rem',
                  color: 'var(--color-gray-600)',
                  margin: 0,
                }}>
                  Gestão de colaboradores
                </p>
              </div>
            </div>
            <div style={{
              display: 'flex',
              alignItems: 'center',
              gap: 'var(--space-3)',
            }}>
              <span style={{
                fontSize: '0.875rem',
                color: 'var(--color-gray-700)',
              }}>
                Olá, {user?.first_name} {user?.last_name}
              </span>
              <button
                onClick={logout}
                style={{
                  backgroundColor: 'var(--color-error)',
                  color: 'var(--text-on-signal)',
                  padding: 'var(--space-2) var(--space-3)',
                  borderRadius: 'var(--radius)',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  border: 'none',
                  cursor: 'pointer',
                  transition: 'background-color var(--transition-fast)',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = 'var(--brand-signal-hover)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = 'var(--color-error)';
                }}
              >
                Terminar sessão
              </button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main style={{
        maxWidth: '1200px',
        margin: '0 auto',
        padding: 'var(--space-6) var(--space-4)',
      }}>
        {/* Controls */}
        <div style={{
          marginBottom: 'var(--space-6)',
          display: 'flex',
          flexDirection: 'column',
          gap: 'var(--space-4)',
        }}>
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            flexWrap: 'wrap',
            gap: 'var(--space-4)',
          }}>
            <div style={{
              position: 'relative',
              maxWidth: '320px',
              flex: '1',
            }}>
              <input
                type="text"
                placeholder="Pesquisar colaboradores…"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                style={{
                  width: '100%',
                  paddingLeft: '2.25rem',
                  paddingRight: 'var(--space-3)',
                  paddingTop: 'var(--space-2)',
                  paddingBottom: 'var(--space-2)',
                  border: '1px solid var(--color-gray-300)',
                  borderRadius: 'var(--radius)',
                  fontSize: '0.875rem',
                  outline: 'none',
                  transition: 'border-color var(--transition-fast)',
                }}
                onFocus={(e) => {
                  e.target.style.borderColor = 'var(--color-primary)';
                  e.target.style.boxShadow = '0 0 0 3px rgba(37, 99, 235, 0.1)';
                }}
                onBlur={(e) => {
                  e.target.style.borderColor = 'var(--color-gray-300)';
                  e.target.style.boxShadow = 'none';
                }}
              />
              <svg 
                style={{
                  position: 'absolute',
                  left: 'var(--space-2)',
                  top: '50%',
                  transform: 'translateY(-50%)',
                  width: '16px',
                  height: '16px',
                  color: 'var(--color-gray-400)',
                }} 
                fill="none" 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <button
              onClick={() => setShowCreateModal(true)}
              style={{
                backgroundColor: 'var(--color-primary)',
                color: 'var(--text-on-primary)',
                padding: 'var(--space-2) var(--space-4)',
                borderRadius: 'var(--radius)',
                fontSize: '0.875rem',
                fontWeight: '500',
                border: 'none',
                cursor: 'pointer',
                transition: 'background-color var(--transition-fast)',
                display: 'flex',
                alignItems: 'center',
                gap: 'var(--space-2)',
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = 'var(--color-primary-hover)';
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = 'var(--color-primary)';
              }}
            >
              <svg 
                style={{ width: '16px', height: '16px' }} 
                fill="none" 
                viewBox="0 0 24 24" 
                stroke="currentColor"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Adicionar colaborador
            </button>
          </div>
        </div>

        {/* Error Message */}
        {error && (
          <div style={{
            marginBottom: 'var(--space-4)',
            backgroundColor: 'var(--chip-danger-bg)',
            border: '1px solid var(--chip-danger-border)',
            color: 'var(--color-error)',
            padding: 'var(--space-3)',
            borderRadius: 'var(--radius)',
            display: 'flex',
            alignItems: 'center',
            gap: 'var(--space-2)',
          }}>
            <svg 
              style={{ width: '16px', height: '16px', flexShrink: 0 }} 
              viewBox="0 0 20 20" 
              fill="currentColor"
            >
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <span style={{ fontSize: '0.875rem' }}>{error}</span>
          </div>
        )}

        {/* Employee Table Container */}
        <div style={{
          backgroundColor: 'var(--surface-header)',
          borderRadius: 'var(--radius-lg)',
          boxShadow: 'var(--shadow-md)',
          overflow: 'hidden',
        }}>
          {/* Table Header */}
          <div style={{
            padding: 'var(--space-4) var(--space-6)',
            borderBottom: '1px solid var(--color-gray-200)',
          }}>
            <h3 style={{
              fontSize: '1.125rem',
              fontWeight: '600',
              color: 'var(--color-gray-900)',
              margin: 0,
              marginBottom: 'var(--space-1)',
            }}>
              Colaboradores ({totalEmployees})
            </h3>
            <p style={{
              fontSize: '0.875rem',
              color: 'var(--color-gray-600)',
              margin: 0,
            }}>
              Gerir a equipa da oficina
            </p>
          </div>
          
          {/* Table */}
          <div style={{ overflowX: 'auto' }}>
            <table style={{
              width: '100%',
              borderCollapse: 'collapse',
            }}>
              <thead style={{
                backgroundColor: 'var(--surface-page)',
              }}>
                <tr>
                  <th style={{
                    padding: 'var(--space-3) var(--space-6)',
                    textAlign: 'left',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: 'var(--color-gray-500)',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}>
                    Nome
                  </th>
                  <th style={{
                    padding: 'var(--space-3) var(--space-6)',
                    textAlign: 'left',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: 'var(--color-gray-500)',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}>
                    E-mail
                  </th>
                  <th style={{
                    padding: 'var(--space-3) var(--space-6)',
                    textAlign: 'left',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: 'var(--color-gray-500)',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}>
                    Departamento
                  </th>
                  <th style={{
                    padding: 'var(--space-3) var(--space-6)',
                    textAlign: 'left',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: 'var(--color-gray-500)',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}>
                    Cargo
                  </th>
                  <th style={{
                    padding: 'var(--space-3) var(--space-6)',
                    textAlign: 'left',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: 'var(--color-gray-500)',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}>
                    Salário
                  </th>
                  <th style={{
                    padding: 'var(--space-3) var(--space-6)',
                    textAlign: 'left',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: 'var(--color-gray-500)',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}>
                    Estado
                  </th>
                  <th style={{
                    padding: 'var(--space-3) var(--space-6)',
                    textAlign: 'left',
                    fontSize: '0.75rem',
                    fontWeight: '500',
                    color: 'var(--color-gray-500)',
                    textTransform: 'uppercase',
                    letterSpacing: '0.05em',
                  }}>
                    Ações
                  </th>
                </tr>
              </thead>
              <tbody style={{
                backgroundColor: 'var(--surface-header)',
              }}>
                {filteredEmployees.length === 0 ? (
                  <tr>
                    <td 
                      colSpan={7} 
                      style={{
                        padding: 'var(--space-8) var(--space-6)',
                        textAlign: 'center',
                        color: 'var(--color-gray-500)',
                        fontSize: '0.875rem',
                      }}
                    >
                      {searchTerm ? 'Nenhum colaborador corresponde à pesquisa.' : 'Sem colaboradores. Adicione o primeiro!'}
                    </td>
                  </tr>
                ) : (
                  filteredEmployees.map((employee, index) => (
                    <tr 
                      key={employee.id}
                      style={{
                        borderTop: index > 0 ? '1px solid var(--color-gray-200)' : 'none',
                        transition: 'background-color var(--transition-fast)',
                      }}
                      onMouseEnter={(e) => {
                        e.currentTarget.style.backgroundColor = 'var(--surface-muted)';
                      }}
                      onMouseLeave={(e) => {
                        e.currentTarget.style.backgroundColor = 'var(--surface-header)';
                      }}
                    >
                      <td style={{
                        padding: 'var(--space-4) var(--space-6)',
                        whiteSpace: 'nowrap',
                      }}>
                        <div style={{
                          display: 'flex',
                          alignItems: 'center',
                          gap: 'var(--space-3)',
                        }}>
                          <div style={{
                            width: '32px',
                            height: '32px',
                            backgroundColor: 'var(--color-primary)',
                            borderRadius: '50%',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            flexShrink: 0,
                          }}>
                            <span style={{
                              fontSize: '0.75rem',
                              fontWeight: '500',
                              color: 'var(--text-on-primary)',
                            }}>
                              {employee.first_name[0]}{employee.last_name[0]}
                            </span>
                          </div>
                          <div>
                            <div style={{
                              fontSize: '0.875rem',
                              fontWeight: '500',
                              color: 'var(--color-gray-900)',
                            }}>
                              {employee.first_name} {employee.last_name}
                            </div>
                            <div style={{
                              fontSize: '0.75rem',
                              color: 'var(--color-gray-500)',
                            }}>
                              {employee.phone || 'No phone'}
                            </div>
                          </div>
                        </div>
                      </td>
                      <td style={{
                        padding: 'var(--space-4) var(--space-6)',
                        whiteSpace: 'nowrap',
                        fontSize: '0.875rem',
                        color: 'var(--color-gray-900)',
                      }}>
                        {employee.email}
                      </td>
                      <td style={{
                        padding: 'var(--space-4) var(--space-6)',
                        whiteSpace: 'nowrap',
                        fontSize: '0.875rem',
                        color: 'var(--color-gray-900)',
                      }}>
                        {employee.department}
                      </td>
                      <td style={{
                        padding: 'var(--space-4) var(--space-6)',
                        whiteSpace: 'nowrap',
                        fontSize: '0.875rem',
                        color: 'var(--color-gray-900)',
                      }}>
                        {employee.position}
                      </td>
                      <td style={{
                        padding: 'var(--space-4) var(--space-6)',
                        whiteSpace: 'nowrap',
                        fontSize: '0.875rem',
                        color: 'var(--color-gray-900)',
                      }}>
                        ${employee.salary.toLocaleString()}
                      </td>
                      <td style={{
                        padding: 'var(--space-4) var(--space-6)',
                        whiteSpace: 'nowrap',
                      }}>
                        <span style={{
                          display: 'inline-flex',
                          padding: 'var(--space-1) var(--space-2)',
                          fontSize: '0.75rem',
                          fontWeight: '500',
                          borderRadius: 'var(--radius)',
                          backgroundColor: employee.is_active ? 'var(--chip-success-bg)' : 'var(--chip-danger-bg)',
                          color: employee.is_active ? 'var(--chip-success-fg)' : 'var(--chip-danger-fg)',
                        }}>
                          {employee.is_active ? 'Ativo' : 'Inativo'}
                        </span>
                      </td>
                      <td style={{
                        padding: 'var(--space-4) var(--space-6)',
                        whiteSpace: 'nowrap',
                      }}>
                        <div style={{
                          display: 'flex',
                          gap: 'var(--space-2)',
                        }}>
                          <button
                            onClick={() => setEditingEmployee(employee)}
                            style={{
                              color: 'var(--color-primary)',
                              backgroundColor: 'transparent',
                              border: 'none',
                              cursor: 'pointer',
                              fontSize: '0.875rem',
                              fontWeight: '500',
                              textDecoration: 'none',
                              transition: 'color var(--transition-fast)',
                            }}
                            onMouseEnter={(e) => {
                              e.currentTarget.style.color = 'var(--color-primary-hover)';
                            }}
                            onMouseLeave={(e) => {
                              e.currentTarget.style.color = 'var(--color-primary)';
                            }}
                          >
                            Editar
                          </button>
                          <button
                            onClick={() => handleDelete(employee.id)}
                            style={{
                              color: 'var(--color-error)',
                              backgroundColor: 'transparent',
                              border: 'none',
                              cursor: 'pointer',
                              fontSize: '0.875rem',
                              fontWeight: '500',
                              textDecoration: 'none',
                              transition: 'color var(--transition-fast)',
                            }}
                            onMouseEnter={(e) => {
                              e.currentTarget.style.color = 'var(--brand-signal-hover)';
                            }}
                            onMouseLeave={(e) => {
                              e.currentTarget.style.color = 'var(--color-error)';
                            }}
                          >
                            Eliminar
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div style={{
              backgroundColor: 'var(--surface-header)',
              padding: 'var(--space-3) var(--space-6)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              borderTop: '1px solid var(--color-gray-200)',
            }}>
              <div>
                <p style={{
                  fontSize: '0.875rem',
                  color: 'var(--color-gray-700)',
                  margin: 0,
                }}>
                  A mostrar{' '}
                  <span style={{ fontWeight: '500' }}>{(currentPage - 1) * employeesPerPage + 1}</span>
                  {' '}a{' '}
                  <span style={{ fontWeight: '500' }}>
                    {Math.min(currentPage * employeesPerPage, totalEmployees)}
                  </span>
                  {' '}de{' '}
                  <span style={{ fontWeight: '500' }}>{totalEmployees}</span>
                  {' '}resultados
                </p>
              </div>
              <div style={{
                display: 'flex',
                gap: 'var(--space-1)',
              }}>
                <button
                  onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                  disabled={currentPage === 1}
                  style={{
                    padding: 'var(--space-2) var(--space-3)',
                    border: '1px solid var(--color-gray-300)',
                    backgroundColor: 'var(--surface-header)',
                    fontSize: '0.875rem',
                    fontWeight: '500',
                    color: currentPage === 1 ? 'var(--color-gray-400)' : 'var(--color-gray-700)',
                    cursor: currentPage === 1 ? 'not-allowed' : 'pointer',
                    borderRadius: 'var(--radius)',
                    transition: 'background-color var(--transition-fast)',
                  }}
                  onMouseEnter={(e) => {
                    if (currentPage !== 1) {
                      e.currentTarget.style.backgroundColor = 'var(--surface-muted)';
                    }
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'var(--surface-header)';
                  }}
                >
                  Anterior
                </button>
                {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
                  const pageNum = i + 1;
                  return (
                    <button
                      key={pageNum}
                      onClick={() => setCurrentPage(pageNum)}
                      style={{
                        padding: 'var(--space-2) var(--space-3)',
                        border: '1px solid var(--color-gray-300)',
                        backgroundColor:
                          pageNum === currentPage ? 'var(--color-primary)' : 'var(--surface-panel)',
                        fontSize: '0.875rem',
                        fontWeight: '500',
                        color:
                          pageNum === currentPage ? 'var(--text-on-primary)' : 'var(--color-gray-700)',
                        cursor: 'pointer',
                        borderRadius: 'var(--radius)',
                        transition: 'all var(--transition-fast)',
                      }}
                      onMouseEnter={(e) => {
                        if (pageNum !== currentPage) {
                          e.currentTarget.style.backgroundColor = 'var(--surface-muted)';
                        }
                      }}
                      onMouseLeave={(e) => {
                        if (pageNum !== currentPage) {
                          e.currentTarget.style.backgroundColor = 'var(--surface-panel)';
                        }
                      }}
                    >
                      {pageNum}
                    </button>
                  );
                })}
                <button
                  onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                  disabled={currentPage === totalPages}
                  style={{
                    padding: 'var(--space-2) var(--space-3)',
                    border: '1px solid var(--color-gray-300)',
                    backgroundColor: 'var(--surface-header)',
                    fontSize: '0.875rem',
                    fontWeight: '500',
                    color: currentPage === totalPages ? 'var(--color-gray-400)' : 'var(--color-gray-700)',
                    cursor: currentPage === totalPages ? 'not-allowed' : 'pointer',
                    borderRadius: 'var(--radius)',
                    transition: 'background-color var(--transition-fast)',
                  }}
                  onMouseEnter={(e) => {
                    if (currentPage !== totalPages) {
                      e.currentTarget.style.backgroundColor = 'var(--surface-muted)';
                    }
                  }}
                  onMouseLeave={(e) => {
                    e.currentTarget.style.backgroundColor = 'var(--surface-header)';
                  }}
                >
                  Seguinte
                </button>
              </div>
            </div>
          )}
        </div>
      </main>

      {/* Create/Edit Employee Modal */}
      {(showCreateModal || editingEmployee) && (
        <EmployeeModal
          employee={editingEmployee}
          onClose={() => {
            setShowCreateModal(false);
            setEditingEmployee(null);
          }}
          onSuccess={() => {
            fetchEmployees();
            setShowCreateModal(false);
            setEditingEmployee(null);
          }}
        />
      )}

      <style jsx>{`
        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
}

// Employee Modal Component with consistent styling
interface EmployeeModalProps {
  employee?: Employee | null;
  onClose: () => void;
  onSuccess: () => void;
}

function EmployeeModal({ employee, onClose, onSuccess }: EmployeeModalProps) {
  const [formData, setFormData] = useState<CreateEmployeeRequest>({
    first_name: employee?.first_name || '',
    last_name: employee?.last_name || '',
    email: employee?.email || '',
    phone: employee?.phone || '',
    department: employee?.department || '',
    position: employee?.position || '',
    hire_date: employee?.hire_date ? employee.hire_date.split('T')[0] : new Date().toISOString().split('T')[0],
    salary: employee?.salary || 0,
  });
  const [errors, setErrors] = useState<{ [key: string]: string }>({});
  const [isLoading, setIsLoading] = useState(false);

  const validateForm = () => {
    const newErrors: { [key: string]: string } = {};

    if (!formData.first_name) newErrors.first_name = 'O nome é obrigatório';
    if (!formData.last_name) newErrors.last_name = 'O apelido é obrigatório';
    if (!formData.email) newErrors.email = 'O e-mail é obrigatório';
    else if (!/\S+@\S+\.\S+/.test(formData.email)) newErrors.email = 'E-mail inválido';
    if (!formData.department) newErrors.department = 'O departamento é obrigatório';
    if (!formData.position) newErrors.position = 'O cargo é obrigatório';
    if (!formData.hire_date) newErrors.hire_date = 'A data de admissão é obrigatória';
    if (!formData.salary || formData.salary <= 0) newErrors.salary = 'O salário tem de ser superior a 0';

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return;

    setIsLoading(true);
    setErrors({});

    try {
      let result;
      if (employee) {
        result = await apiClient.updateEmployee(employee.id, formData);
      } else {
        result = await apiClient.createEmployee(formData);
      }

      if (result.data && !result.error) {
        onSuccess();
      } else {
        setErrors({ general: result.error?.message || 'Operação falhou' });
      }
    } catch (error) {
      setErrors({ general: 'Erro de rede. Tente novamente.' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'salary' ? parseFloat(value) || 0 : value
    }));
    
    if (errors[name]) {
      setErrors(prev => ({ ...prev, [name]: '' }));
    }
  };

  return (
    <div style={{
      position: 'fixed',
      inset: 0,
      backgroundColor: 'rgba(0, 0, 0, 0.5)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: 'var(--space-4)',
      zIndex: 50,
    }}>
      <div style={{
        backgroundColor: 'var(--surface-header)',
        borderRadius: 'var(--radius-lg)',
        boxShadow: 'var(--shadow-lg)',
        width: '100%',
        maxWidth: '500px',
        maxHeight: '90vh',
        overflow: 'hidden',
      }}>
        {/* Modal Header */}
        <div style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: 'var(--space-4) var(--space-6)',
          borderBottom: '1px solid var(--color-gray-200)',
        }}>
          <h3 style={{
            fontSize: '1.125rem',
            fontWeight: '600',
            color: 'var(--color-gray-900)',
            margin: 0,
          }}>
            {employee ? 'Editar colaborador' : 'Novo colaborador'}
          </h3>
          <button
            onClick={onClose}
            style={{
              color: 'var(--color-gray-400)',
              backgroundColor: 'transparent',
              border: 'none',
              cursor: 'pointer',
              padding: 'var(--space-1)',
              borderRadius: 'var(--radius)',
              transition: 'color var(--transition-fast)',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.color = 'var(--color-gray-600)';
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.color = 'var(--color-gray-400)';
            }}
          >
            <svg style={{ width: '20px', height: '20px' }} fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Modal Body */}
        <div style={{
          padding: 'var(--space-6)',
          maxHeight: '70vh',
          overflowY: 'auto',
        }}>
          <form onSubmit={handleSubmit}>
            {errors.general && (
              <div style={{
                backgroundColor: 'var(--chip-danger-bg)',
                border: '1px solid var(--chip-danger-border)',
                color: 'var(--color-error)',
                padding: 'var(--space-3)',
                borderRadius: 'var(--radius)',
                marginBottom: 'var(--space-4)',
                fontSize: '0.875rem',
              }}>
                {errors.general}
              </div>
            )}

            <div style={{
              display: 'grid',
              gridTemplateColumns: '1fr 1fr',
              gap: 'var(--space-4)',
              marginBottom: 'var(--space-4)',
            }}>
              <div>
                <label style={{
                  display: 'block',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--color-gray-700)',
                  marginBottom: 'var(--space-1)',
                }}>
                  Nome
                </label>
                <input
                  type="text"
                  name="first_name"
                  value={formData.first_name}
                  onChange={handleChange}
                  style={{
                    width: '100%',
                    padding: 'var(--space-2)',
                    border: `1px solid ${errors.first_name ? 'var(--chip-danger-border)' : 'var(--color-gray-300)'}`,
                    borderRadius: 'var(--radius)',
                    fontSize: '0.875rem',
                    outline: 'none',
                    transition: 'border-color var(--transition-fast)',
                  }}
                  onFocus={(e) => {
                    e.target.style.borderColor = 'var(--color-primary)';
                    e.target.style.boxShadow = '0 0 0 3px rgba(37, 99, 235, 0.1)';
                  }}
                  onBlur={(e) => {
                    e.target.style.borderColor = errors.first_name ? 'var(--chip-danger-border)' : 'var(--color-gray-300)';
                    e.target.style.boxShadow = 'none';
                  }}
                />
                {errors.first_name && (
                  <p style={{
                    marginTop: 'var(--space-1)',
                    fontSize: '0.75rem',
                    color: 'var(--color-error)',
                  }}>
                    {errors.first_name}
                  </p>
                )}
              </div>

              <div>
                <label style={{
                  display: 'block',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--color-gray-700)',
                  marginBottom: 'var(--space-1)',
                }}>
                  Apelido
                </label>
                <input
                  type="text"
                  name="last_name"
                  value={formData.last_name}
                  onChange={handleChange}
                  style={{
                    width: '100%',
                    padding: 'var(--space-2)',
                    border: `1px solid ${errors.last_name ? 'var(--chip-danger-border)' : 'var(--color-gray-300)'}`,
                    borderRadius: 'var(--radius)',
                    fontSize: '0.875rem',
                    outline: 'none',
                    transition: 'border-color var(--transition-fast)',
                  }}
                  onFocus={(e) => {
                    e.target.style.borderColor = 'var(--color-primary)';
                    e.target.style.boxShadow = '0 0 0 3px rgba(37, 99, 235, 0.1)';
                  }}
                  onBlur={(e) => {
                    e.target.style.borderColor = errors.last_name ? 'var(--chip-danger-border)' : 'var(--color-gray-300)';
                    e.target.style.boxShadow = 'none';
                  }}
                />
                {errors.last_name && (
                  <p style={{
                    marginTop: 'var(--space-1)',
                    fontSize: '0.75rem',
                    color: 'var(--color-error)',
                  }}>
                    {errors.last_name}
                  </p>
                )}
              </div>
            </div>

            <div style={{ marginBottom: 'var(--space-4)' }}>
              <label style={{
                display: 'block',
                fontSize: '0.875rem',
                fontWeight: '500',
                color: 'var(--color-gray-700)',
                marginBottom: 'var(--space-1)',
              }}>
                E-mail
              </label>
              <input
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
                style={{
                  width: '100%',
                  padding: 'var(--space-2)',
                  border: `1px solid ${errors.email ? 'var(--chip-danger-border)' : 'var(--color-gray-300)'}`,
                  borderRadius: 'var(--radius)',
                  fontSize: '0.875rem',
                  outline: 'none',
                  transition: 'border-color var(--transition-fast)',
                }}
                onFocus={(e) => {
                  e.target.style.borderColor = 'var(--color-primary)';
                  e.target.style.boxShadow = '0 0 0 3px rgba(37, 99, 235, 0.1)';
                }}
                onBlur={(e) => {
                  e.target.style.borderColor = errors.email ? 'var(--chip-danger-border)' : 'var(--color-gray-300)';
                  e.target.style.boxShadow = 'none';
                }}
              />
              {errors.email && (
                <p style={{
                  marginTop: 'var(--space-1)',
                  fontSize: '0.75rem',
                  color: 'var(--color-error)',
                }}>
                  {errors.email}
                </p>
              )}
            </div>

            <div style={{
              display: 'grid',
              gridTemplateColumns: '1fr 1fr',
              gap: 'var(--space-4)',
              marginBottom: 'var(--space-4)',
            }}>
              <div>
                <label style={{
                  display: 'block',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--color-gray-700)',
                  marginBottom: 'var(--space-1)',
                }}>
                  Departamento
                </label>
                <select
                  name="department"
                  value={formData.department}
                  onChange={handleChange}
                  style={{
                    width: '100%',
                    padding: 'var(--space-2)',
                    border: `1px solid ${errors.department ? 'var(--chip-danger-border)' : 'var(--color-gray-300)'}`,
                    borderRadius: 'var(--radius)',
                    fontSize: '0.875rem',
                    outline: 'none',
                    transition: 'border-color var(--transition-fast)',
                  }}
                >
                  <option value="">Selecione o departamento</option>
                  <option value="Mechanical">Mecânica</option>
                  <option value="Electrical">Elétrica</option>
                  <option value="Body Work">Chapa e pintura</option>
                  <option value="Administration">Administração</option>
                  <option value="Sales">Vendas</option>
                </select>
                {errors.department && (
                  <p style={{
                    marginTop: 'var(--space-1)',
                    fontSize: '0.75rem',
                    color: 'var(--color-error)',
                  }}>
                    {errors.department}
                  </p>
                )}
              </div>

              <div>
                <label style={{
                  display: 'block',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--color-gray-700)',
                  marginBottom: 'var(--space-1)',
                }}>
                  Cargo
                </label>
                <input
                  type="text"
                  name="position"
                  value={formData.position}
                  onChange={handleChange}
                  style={{
                    width: '100%',
                    padding: 'var(--space-2)',
                    border: `1px solid ${errors.position ? 'var(--chip-danger-border)' : 'var(--color-gray-300)'}`,
                    borderRadius: 'var(--radius)',
                    fontSize: '0.875rem',
                    outline: 'none',
                    transition: 'border-color var(--transition-fast)',
                  }}
                  onFocus={(e) => {
                    e.target.style.borderColor = 'var(--color-primary)';
                    e.target.style.boxShadow = '0 0 0 3px rgba(37, 99, 235, 0.1)';
                  }}
                  onBlur={(e) => {
                    e.target.style.borderColor = errors.position ? 'var(--chip-danger-border)' : 'var(--color-gray-300)';
                    e.target.style.boxShadow = 'none';
                  }}
                />
                {errors.position && (
                  <p style={{
                    marginTop: 'var(--space-1)',
                    fontSize: '0.75rem',
                    color: 'var(--color-error)',
                  }}>
                    {errors.position}
                  </p>
                )}
              </div>
            </div>

            <div style={{
              display: 'grid',
              gridTemplateColumns: '1fr 1fr',
              gap: 'var(--space-4)',
              marginBottom: 'var(--space-4)',
            }}>
              <div>
                <label style={{
                  display: 'block',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--color-gray-700)',
                  marginBottom: 'var(--space-1)',
                }}>
                  Data de admissão
                </label>
                <input
                  type="date"
                  name="hire_date"
                  value={formData.hire_date}
                  onChange={handleChange}
                  style={{
                    width: '100%',
                    padding: 'var(--space-2)',
                    border: `1px solid ${errors.hire_date ? 'var(--chip-danger-border)' : 'var(--color-gray-300)'}`,
                    borderRadius: 'var(--radius)',
                    fontSize: '0.875rem',
                    outline: 'none',
                    transition: 'border-color var(--transition-fast)',
                  }}
                  onFocus={(e) => {
                    e.target.style.borderColor = 'var(--color-primary)';
                    e.target.style.boxShadow = '0 0 0 3px rgba(37, 99, 235, 0.1)';
                  }}
                  onBlur={(e) => {
                    e.target.style.borderColor = errors.hire_date ? 'var(--chip-danger-border)' : 'var(--color-gray-300)';
                    e.target.style.boxShadow = 'none';
                  }}
                />
                {errors.hire_date && (
                  <p style={{
                    marginTop: 'var(--space-1)',
                    fontSize: '0.75rem',
                    color: 'var(--color-error)',
                  }}>
                    {errors.hire_date}
                  </p>
                )}
              </div>

              <div>
                <label style={{
                  display: 'block',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--color-gray-700)',
                  marginBottom: 'var(--space-1)',
                }}>
                  Salário
                </label>
                <input
                  type="number"
                  name="salary"
                  value={formData.salary}
                  onChange={handleChange}
                  min="0"
                  step="0.01"
                  style={{
                    width: '100%',
                    padding: 'var(--space-2)',
                    border: `1px solid ${errors.salary ? 'var(--chip-danger-border)' : 'var(--color-gray-300)'}`,
                    borderRadius: 'var(--radius)',
                    fontSize: '0.875rem',
                    outline: 'none',
                    transition: 'border-color var(--transition-fast)',
                  }}
                  onFocus={(e) => {
                    e.target.style.borderColor = 'var(--color-primary)';
                    e.target.style.boxShadow = '0 0 0 3px rgba(37, 99, 235, 0.1)';
                  }}
                  onBlur={(e) => {
                    e.target.style.borderColor = errors.salary ? 'var(--chip-danger-border)' : 'var(--color-gray-300)';
                    e.target.style.boxShadow = 'none';
                  }}
                />
                {errors.salary && (
                  <p style={{
                    marginTop: 'var(--space-1)',
                    fontSize: '0.75rem',
                    color: 'var(--color-error)',
                  }}>
                    {errors.salary}
                  </p>
                )}
              </div>
            </div>

            <div style={{
              display: 'flex',
              justifyContent: 'flex-end',
              gap: 'var(--space-3)',
              paddingTop: 'var(--space-4)',
              borderTop: '1px solid var(--color-gray-200)',
            }}>
              <button
                type="button"
                onClick={onClose}
                style={{
                  padding: 'var(--space-2) var(--space-4)',
                  border: '1px solid var(--color-gray-300)',
                  borderRadius: 'var(--radius)',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--color-gray-700)',
                  backgroundColor: 'var(--surface-header)',
                  cursor: 'pointer',
                  transition: 'background-color var(--transition-fast)',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = 'var(--surface-muted)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = 'var(--surface-header)';
                }}
              >
                Cancelar
              </button>
              <button
                type="submit"
                disabled={isLoading}
                style={{
                  padding: 'var(--space-2) var(--space-4)',
                  border: 'none',
                  borderRadius: 'var(--radius)',
                  fontSize: '0.875rem',
                  fontWeight: '500',
                  color: 'var(--text-on-primary)',
                  backgroundColor: isLoading ? 'var(--color-gray-400)' : 'var(--color-primary)',
                  cursor: isLoading ? 'not-allowed' : 'pointer',
                  transition: 'background-color var(--transition-fast)',
                }}
                onMouseEnter={(e) => {
                  if (!isLoading) {
                    e.currentTarget.style.backgroundColor = 'var(--color-primary-hover)';
                  }
                }}
                onMouseLeave={(e) => {
                  if (!isLoading) {
                    e.currentTarget.style.backgroundColor = 'var(--color-primary)';
                  }
                }}
              >
                {isLoading ? 'A guardar…' : (employee ? 'Atualizar' : 'Criar')}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}