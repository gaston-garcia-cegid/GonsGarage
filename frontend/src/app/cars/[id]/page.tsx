'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useAuth, useCars } from '@/stores';
import { apiClient, Repair } from '@/lib/api';
import AppShell from '@/components/layout/AppShell';
import { useAuthHydrationReady } from '@/hooks/useAuthHydrationReady';
import styles from './car-details.module.css';

function repairStatusPt(status: string): string {
  const map: Record<string, string> = {
    pending: 'Pendente',
    in_progress: 'Em curso',
    completed: 'Concluída',
    cancelled: 'Cancelada',
  };
  return map[status] ?? status.replace(/_/g, ' ');
}

export default function CarDetailsPage() {
  const [repairs, setRepairs] = useState<Repair[]>([]);
  const [repairsLoading, setRepairsLoading] = useState(true);

  const { user, logout } = useAuth();
  const authHydrated = useAuthHydrationReady();
  const { selectedCar: car, isLoading: loading, error, fetchCarById } = useCars();
  const router = useRouter();
  const params = useParams();
  const carId = params.id as string;

  useEffect(() => {
    if (carId === 'new') {
      router.replace('/cars?addCar=1');
    }
  }, [carId, router]);

  const fetchCarRepairs = useCallback(async () => {
    try {
      setRepairsLoading(true);

      const { data, error: apiError } = await apiClient.getRepairs(carId);
      
      if (data && !apiError) {
        setRepairs(data);
      } else {
        console.error('Failed to fetch repairs:', apiError);
      }
    } catch (err) {
      console.error('Network error fetching repairs:', err);
    } finally {
      setRepairsLoading(false);
    }
  }, [carId]);

  useEffect(() => {
    if (!authHydrated) return;
    if (!user) {
      router.replace('/auth/login');
      return;
    }
    if (!carId || carId === 'new') return;
    fetchCarById(carId);
    fetchCarRepairs();
  }, [authHydrated, user, router, carId, fetchCarById, fetchCarRepairs]);

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-PT', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('pt-PT', {
      style: 'currency',
      currency: 'EUR',
    }).format(amount);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'warning';
      case 'in_progress': return 'info';
      case 'completed': return 'success';
      case 'cancelled': return 'error';
      default: return 'default';
    }
  };

  if (!authHydrated || !user) {
    return (
      <div className="loadingScreen">
        <div className="spinnerLg" aria-hidden />
      </div>
    );
  }

  if (carId === 'new') {
    return (
      <div className="loadingScreen">
        <div className="spinnerLg" aria-hidden />
      </div>
    );
  }

  if (loading) {
    return (
      <AppShell user={user} subtitle="Detalhe do automóvel" activeNav="cars" onLogout={logout}>
        <div className={styles.loadingContainer}>
          <div className={styles.spinner}></div>
          <span>A carregar detalhes…</span>
        </div>
      </AppShell>
    );
  }

  if (error || !car) {
    return (
      <AppShell user={user} subtitle="Detalhe do automóvel" activeNav="cars" onLogout={logout}>
        <div className={styles.errorContainer}>
          <div className={styles.errorContent}>
            <div className={styles.errorIcon}>⚠️</div>
            <h2>Automóvel não encontrado</h2>
            <p>{error || 'O automóvel pedido não existe ou não está disponível.'}</p>
            <button type="button" onClick={() => router.push('/cars')} className={styles.backButton}>
              Voltar aos automóveis
            </button>
          </div>
        </div>
      </AppShell>
    );
  }

  return (
    <AppShell user={user} subtitle="Detalhe do automóvel" activeNav="cars" onLogout={logout}>
      <div className={styles.breadcrumbs}>
        <button onClick={() => router.push('/cars')} className={styles.breadcrumbLink}>
          Os meus automóveis
        </button>
        <span className={styles.breadcrumbSeparator}>›</span>
        <span className={styles.breadcrumbCurrent}>
          {car.year} {car.make} {car.model}
        </span>
      </div>

      <div className={styles.main}>
        {/* Car Info Card */}
        <div className={styles.carInfoCard}>
          <div className={styles.carInfoHeader}>
            <div className={styles.carIconLarge}>🚗</div>
            <div className={styles.carTitle}>
              <h2>{car.year} {car.make} {car.model}</h2>
              <p className={styles.licensePlate}>{car.licensePlate}</p>
            </div>
            <div className={styles.carActions}>
              <button 
                onClick={() => router.push(`/appointments?schedule=1&carId=${encodeURIComponent(car.id)}`)}
                className={styles.scheduleServiceButton}
              >
                <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                Marcar serviço
              </button>
            </div>
          </div>

          <div className={styles.carInfoGrid}>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>Cor</span>
              <span className={styles.infoValue}>{car.color}</span>
            </div>
            {car.vin && (
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>VIN</span>
                <span className={styles.infoValue}>{car.vin}</span>
              </div>
            )}
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>Registado em</span>
              <span className={styles.infoValue}>{formatDate(car.createdAt)}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>Última atualização</span>
              <span className={styles.infoValue}>{formatDate(car.updatedAt)}</span>
            </div>
          </div>
        </div>

        {/* Repairs Section */}
        <div className={styles.repairsSection}>
          <div className={styles.sectionHeader}>
            <h3>Histórico de serviços</h3>
            <div className={styles.repairsStats}>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {repairs.filter(r => r.status === 'completed').length}
                </span>
                <span className={styles.statLabel}>Concluídas</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {repairs.filter(r => r.status === 'in_progress').length}
                </span>
                <span className={styles.statLabel}>Em curso</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {repairs.filter(r => r.status === 'pending').length}
                </span>
                <span className={styles.statLabel}>Pendentes</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {formatCurrency(repairs.reduce((sum, r) => sum + r.cost, 0))}
                </span>
                <span className={styles.statLabel}>Custo total</span>
              </div>
            </div>
          </div>

          {repairsLoading ? (
            <div className={styles.repairsLoading}>
              <div className={styles.spinner}></div>
              <span>A carregar reparações…</span>
            </div>
          ) : repairs.length === 0 ? (
            <div className={styles.emptyRepairs}>
              <div className={styles.emptyIcon}>🔧</div>
              <h4>Sem histórico de serviços</h4>
              <p>Este automóvel ainda não tem intervenções registadas.</p>
              <button 
                onClick={() => router.push(`/appointments?schedule=1&carId=${encodeURIComponent(car.id)}`)}
                className={styles.scheduleFirstServiceButton}
              >
                Marcar primeiro serviço
              </button>
            </div>
          ) : (
            <div className={styles.repairsList}>
              {repairs.map((repair) => (
                <div key={repair.id} className={styles.repairCard}>
                  <div className={styles.repairHeader}>
                    <div className={styles.repairStatus}>
                      <span className={`${styles.statusBadge} ${styles[getStatusColor(repair.status)]}`}>
                        {repairStatusPt(repair.status)}
                      </span>
                    </div>
                    <div className={styles.repairCost}>
                      {formatCurrency(repair.cost)}
                    </div>
                  </div>

                  <div className={styles.repairContent}>
                    <h4 className={styles.repairDescription}>{repair.description}</h4>
                    
                    <div className={styles.repairMeta}>
                      <div className={styles.metaItem}>
                        <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        <span>Criado: {formatDate(repair.created_at)}</span>
                      </div>
                      
                      {repair.started_at && (
                        <div className={styles.metaItem}>
                          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span>Início: {formatDate(repair.started_at)}</span>
                        </div>
                      )}
                      
                      {repair.completed_at && (
                        <div className={styles.metaItem}>
                          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span>Concluído: {formatDate(repair.completed_at)}</span>
                        </div>
                      )}

                      {repair.technician && (
                        <div className={styles.metaItem}>
                          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                          </svg>
                          <span>Técnico: {repair.technician.first_name} {repair.technician.last_name}</span>
                        </div>
                      )}
                    </div>
                  </div>

                  <div className={styles.repairProgress}>
                    <div className={styles.progressBar}>
                      <div 
                        className={`${styles.progressFill} ${styles[getStatusColor(repair.status)]}`}
                        style={{
                          width: repair.status === 'completed' ? '100%' : 
                                repair.status === 'in_progress' ? '60%' : 
                                repair.status === 'pending' ? '20%' : '0%'
                        }}
                      ></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </AppShell>
  );
}