'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '../contexts/AuthContext';
import Image from 'next/image';
import styles from './landing.module.css';

export default function LandingPage() {
  const { isAuthenticated, user, isLoading } = useAuth();
  const router = useRouter();

  // Só redireciona se o usuário ESTIVER autenticado
  useEffect(() => {
    console.log('Landing page - Auth state:', { isLoading, isAuthenticated, userRole: user?.role });
  
    if (!isLoading && isAuthenticated && user) {
      switch (user.role) {
        case 'admin':
        case 'manager':
          console.log('Redirecting to admin dashboard');
          router.push('/admin/dashboard');
          break;
        case 'employee':
        case 'technician':
          console.log('Redirecting to technician dashboard');
          router.push('/technician/dashboard');
          break;
        case 'client':
          console.log('Redirecting to client dashboard');
          router.push('/client/dashboard');
          break;
        default:
          console.log('Unknown role, redirecting to default dashboard');
          router.push('/dashboard');
      }
    }
    // REMOVIDO: router.push('/auth/login') para usuários não autenticados
  }, [isAuthenticated, user, isLoading, router]);

  const handleLogin = () => {
    router.push('/auth/login');
  };

  const handleRegister = () => {
    router.push('/auth/register');
  };

  const services = [
    {
      icon: '🔧',
      title: 'General Maintenance',
      description: 'Regular maintenance to keep your car running smoothly',
      price: 'From $89'
    },
    {
      icon: '🛞',
      title: 'Tire Services',
      description: 'Tire rotation, alignment, and replacement services',
      price: 'From $45'
    },
    {
      icon: '🔋',
      title: 'Battery Service',
      description: 'Battery testing, maintenance, and replacement',
      price: 'From $120'
    },
    {
      icon: '❄️',
      title: 'Air Conditioning',
      description: 'A/C repair, maintenance, and refrigerant refill',
      price: 'From $150'
    },
    {
      icon: '🚗',
      title: 'Engine Diagnostic',
      description: 'Complete engine analysis and troubleshooting',
      price: 'From $95'
    },
    {
      icon: '🛡️',
      title: 'Brake Service',
      description: 'Brake pads, rotors, and complete brake system service',
      price: 'From $180'
    }
  ];

  const features = [
    {
      icon: '👨‍🔧',
      title: 'Expert Technicians',
      description: 'Certified mechanics with years of experience'
    },
    {
      icon: '⚡',
      title: 'Quick Service',
      description: 'Fast and efficient repairs to get you back on the road'
    },
    {
      icon: '💰',
      title: 'Fair Pricing',
      description: 'Transparent pricing with no hidden fees'
    },
    {
      icon: '🛡️',
      title: 'Quality Guarantee',
      description: '90-day warranty on all our services'
    }
  ];

  // Mostra loading apenas enquanto está verificando autenticação
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  // Se o usuário estiver autenticado, não mostra a landing page (vai redirecionar)
  if (isAuthenticated && user) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p>Redirecting to your dashboard...</p>
        </div>
      </div>
    );
  }

  // Mostra a landing page apenas para usuários NÃO autenticados
  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.logoSection}>
            <div className={styles.logo}>
              <Image
                src="/images/LogoGonsGarage.jpg" // ou .jpg/.jpeg dependendo da extensão
                alt="GonsGarage Logo"
                width={32}
                height={32}
                style={{ objectFit: 'contain' }}
              />
            </div>
            <div className={styles.logoText}>
              <h1>GonsGarage</h1>
              <p>Professional Auto Repair</p>
            </div>
          </div>
          
          <nav className={styles.navigation}>
            <a href="#services" className={styles.navLink}>Services</a>
            <a href="#about" className={styles.navLink}>About</a>
            <a href="#contact" className={styles.navLink}>Contact</a>
          </nav>

          <div className={styles.authButtons}>
            <button onClick={handleLogin} className={styles.loginButton}>
              Login
            </button>
            <button onClick={handleRegister} className={styles.registerButton}>
              Sign Up
            </button>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className={styles.hero}>
        <div className={styles.heroBackground}>
          {/* Placeholder para a imagem do banner */}
          <div className={styles.bannerImage}>
            <Image
              src="/images/BannerGonsGarage.jpg"
              alt="GonsGarage Workshop"
              fill
              style={{ objectFit: 'cover' }}
              priority
              onError={(e) => {
                // Fallback se a imagem não carregar
                e.currentTarget.style.display = 'none';
              }}
            />
          </div>
          <div className={styles.heroOverlay} />
        </div>
        
        <div className={styles.heroContent}>
          <div className={styles.heroText}>
            <h2>Your Trusted Auto Repair Shop</h2>
            <p>
              Expert mechanics, quality service, and fair prices. 
              We keep your vehicle running at its best with professional care and attention.
            </p>
            <div className={styles.heroStats}>
              <div className={styles.stat}>
                <span className={styles.statNumber}>15+</span>
                <span className={styles.statLabel}>Years Experience</span>
              </div>
              <div className={styles.stat}>
                <span className={styles.statNumber}>5000+</span>
                <span className={styles.statLabel}>Happy Customers</span>
              </div>
              <div className={styles.stat}>
                <span className={styles.statNumber}>24/7</span>
                <span className={styles.statLabel}>Emergency Service</span>
              </div>
            </div>
            <div className={styles.heroActions}>
              <button onClick={handleRegister} className={styles.ctaButton}>
                Book Service Now
              </button>
              <button onClick={handleLogin} className={styles.secondaryButton}>
                Existing Customer
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* Services Section */}
      <section id="services" className={styles.servicesSection}>
        <div className={styles.sectionContent}>
          <div className={styles.sectionHeader}>
            <h3>Our Services</h3>
            <p>Complete automotive care for all makes and models</p>
          </div>
          
          <div className={styles.servicesGrid}>
            {services.map((service, index) => (
              <div key={index} className={styles.serviceCard}>
                <div className={styles.serviceIcon}>{service.icon}</div>
                <h4>{service.title}</h4>
                <p>{service.description}</p>
                <div className={styles.servicePrice}>{service.price}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className={styles.featuresSection}>
        <div className={styles.sectionContent}>
          <div className={styles.sectionHeader}>
            <h3>Why Choose GonsGarage?</h3>
            <p>We&apos;re committed to providing exceptional automotive service</p>
          </div>
          
          <div className={styles.featuresGrid}>
            {features.map((feature, index) => (
              <div key={index} className={styles.featureCard}>
                <div className={styles.featureIcon}>{feature.icon}</div>
                <h4>{feature.title}</h4>
                <p>{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Contact Section */}
      <section id="contact" className={styles.contactSection}>
        <div className={styles.sectionContent}>
          <div className={styles.contactGrid}>
            <div className={styles.contactInfo}>
              <h3>Visit Our Shop</h3>
              <div className={styles.contactItem}>
                <div className={styles.contactIcon}>📍</div>
                <div>
                  <h4>Address</h4>
                  <p>123 Main Street<br />Downtown, State 12345</p>
                </div>
              </div>
              <div className={styles.contactItem}>
                <div className={styles.contactIcon}>📞</div>
                <div>
                  <h4>Phone</h4>
                  <p>(555) 123-4567</p>
                </div>
              </div>
              <div className={styles.contactItem}>
                <div className={styles.contactIcon}>🕒</div>
                <div>
                  <h4>Hours</h4>
                  <p>Mon-Fri: 8AM-6PM<br />Sat: 8AM-4PM<br />Sun: Closed</p>
                </div>
              </div>
            </div>
            
            <div className={styles.contactCta}>
              <h3>Ready to Get Started?</h3>
              <p>Book your appointment today and experience the GonsGarage difference</p>
              <button onClick={handleRegister} className={styles.ctaButton}>
                Schedule Service
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className={styles.footer}>
        <div className={styles.footerContent}>
          <div className={styles.footerSection}>
            <div className={styles.footerLogo}>
              <div className={styles.logo}>
                <Image
                  src="/images/LogoGonsGarage.jpg"
                  alt="GonsGarage Logo"
                  width={32}
                  height={32}
                  style={{ objectFit: 'contain' }}
                />
              </div>
              <div>
                <h4>GonsGarage</h4>
                <p>Professional Auto Repair</p>
              </div>
            </div>
          </div>
          
          <div className={styles.footerSection}>
            <h4>Services</h4>
            <ul>
              <li>Oil Change</li>
              <li>Brake Service</li>
              <li>Tire Service</li>
              <li>Engine Repair</li>
            </ul>
          </div>
          
          <div className={styles.footerSection}>
            <h4>Company</h4>
            <ul>
              <li>About Us</li>
              <li>Our Team</li>
              <li>Careers</li>
              <li>Contact</li>
            </ul>
          </div>
          
          <div className={styles.footerSection}>
            <h4>Customer</h4>
            <ul>
              <li><button onClick={handleLogin}>Login</button></li>
              <li><button onClick={handleRegister}>Sign Up</button></li>
              <li>Service History</li>
              <li>Support</li>
            </ul>
          </div>
        </div>
        
        <div className={styles.footerBottom}>
          <p>&copy; 2024 GonsGarage. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}
