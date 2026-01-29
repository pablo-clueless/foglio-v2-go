# Foglio API v2

A professional networking and job platform API built with Go, Gin, and PostgreSQL.

## Table of Contents

- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Subscription & Payment Integration](#subscription--payment-integration)
  - [Overview](#overview)
  - [Subscription Tiers](#subscription-tiers)
  - [Payment Flow](#payment-flow)
  - [Frontend Implementation Guide](#frontend-implementation-guide)
  - [Webhook Handling](#webhook-handling)
  - [Error Handling](#error-handling)

---

## Getting Started

```bash
# Install dependencies
go mod download

# Run the server
go run main.go
```

## Environment Variables

```env
# App
GO_ENV=development
PORT=8080
VERSION=api/v2
CLIENT_URL=http://localhost:3000

# Database
POSTGRES_URL=postgresql://user:password@localhost:5432/foglio

# JWT
JWT_SECRET=your-jwt-secret

# Paystack
PAYSTACK_SECRET_KEY=sk_test_xxxxx
PAYSTACK_PUBLIC_KEY=pk_test_xxxxx
PAYSTACK_WEBHOOK_SECRET=whsec_xxxxx
```

---

## Subscription & Payment Integration

### Overview

Foglio uses **Paystack** for payment processing. The subscription system supports:

- Multiple subscription tiers (Free, Basic, Premium, Business)
- Monthly and yearly billing cycles
- Automatic recurring payments
- Subscription upgrades/downgrades
- Cancellation with period-end handling

### Subscription Tiers

#### Get Available Tiers

```http
GET /api/v2/subscriptions
```

**Response:**
```json
{
  "success": true,
  "message": "Subscriptions fetched successfully",
  "data": {
    "data": [
      {
        "id": "uuid-here",
        "name": "Premium Plan",
        "description": "Best for professionals",
        "type": "monthly",
        "tier": "premium",
        "price": 2999.00,
        "currency": "NGN",
        "billing_cycle_days": 30,
        "trial_period_days": 14,
        "features": {
          "max_projects": 50,
          "priority_support": true
        },
        "max_projects": 50,
        "max_skills": 100,
        "max_experiences": 50,
        "is_active": true,
        "is_popular": true
      }
    ],
    "total_items": 4,
    "total_pages": 1,
    "page": 1,
    "limit": 10
  }
}
```

---

### Payment Flow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           SUBSCRIPTION PAYMENT FLOW                          │
└─────────────────────────────────────────────────────────────────────────────┘

    ┌──────────┐         ┌──────────┐         ┌──────────┐         ┌──────────┐
    │ Foglio   │         │  Foglio  │         │ Paystack │         │  Webhook │
    │   App    │         │   API    │         │   API    │         │ Endpoint │
    └────┬─────┘         └────┬─────┘         └────┬─────┘         └────┬─────┘
         │                    │                    │                    │
         │  1. Initialize     │                    │                    │
         │     Payment        │                    │                    │
         │───────────────────>│                    │                    │
         │                    │                    │                    │
         │                    │  2. Create         │                    │
         │                    │     Transaction    │                    │
         │                    │───────────────────>│                    │
         │                    │                    │                    │
         │                    │  3. Return auth_url│                    │
         │                    │<───────────────────│                    │
         │                    │                    │                    │
         │  4. Return         │                    │                    │
         │     auth_url       │                    │                    │
         │<───────────────────│                    │                    │
         │                    │                    │                    │
         │  5. Redirect user  │                    │                    │
         │     to Paystack    │                    │                    │
         │─────────────────────────────────────────>                    │
         │                    │                    │                    │
         │                    │                    │  6. User pays      │
         │                    │                    │     on Paystack    │
         │                    │                    │                    │
         │  7. Redirect back  │                    │                    │
         │     with reference │                    │                    │
         │<─────────────────────────────────────────                    │
         │                    │                    │                    │
         │  8. Verify         │                    │                    │
         │     Payment        │                    │                    │
         │───────────────────>│                    │                    │
         │                    │                    │                    │
         │                    │  9. Verify with    │                    │
         │                    │     Paystack       │                    │
         │                    │───────────────────>│                    │
         │                    │                    │                    │
         │                    │  10. Confirmation  │                    │
         │                    │<───────────────────│                    │
         │                    │                    │                    │
         │                    │  11. Activate      │                    │
         │                    │      Subscription  │                    │
         │                    │                    │                    │
         │  12. Success       │                    │                    │
         │      Response      │                    │                    │
         │<───────────────────│                    │                    │
         │                    │                    │                    │
         │                    │                    │  13. Webhook       │
         │                    │                    │      (backup)      │
         │                    │                    │─────────────────────>
         │                    │                    │                    │
         │                    │                    │                    │
    ┌────┴─────┐         ┌────┴─────┐         ┌────┴─────┐         ┌────┴─────┐
    │ Frontend │         │  Foglio  │         │ Paystack │         │  Webhook │
    │   App    │         │   API    │         │   API    │         │ Endpoint │
    └──────────┘         └──────────┘         └──────────┘         └──────────┘
```

---

### Frontend Implementation Guide

#### Step 1: Display Subscription Plans

Fetch and display available subscription tiers to the user.

```typescript
// TypeScript/React Example

interface SubscriptionTier {
  id: string;
  name: string;
  description: string;
  type: 'monthly' | 'yearly' | 'lifetime';
  tier: 'free' | 'basic' | 'premium' | 'business';
  price: number;
  currency: string;
  billing_cycle_days: number;
  features: Record<string, any>;
  is_popular: boolean;
}

async function fetchSubscriptionTiers(): Promise<SubscriptionTier[]> {
  const response = await fetch('/api/v2/subscriptions', {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  const result = await response.json();
  return result.data.data;
}

// Display in your component
function PricingPage() {
  const [tiers, setTiers] = useState<SubscriptionTier[]>([]);

  useEffect(() => {
    fetchSubscriptionTiers().then(setTiers);
  }, []);

  return (
    <div className="pricing-grid">
      {tiers.map(tier => (
        <PricingCard
          key={tier.id}
          tier={tier}
          onSelect={() => handleSubscribe(tier.id)}
        />
      ))}
    </div>
  );
}
```

---

#### Step 2: Initialize Payment

When user selects a plan, initialize the payment to get the Paystack authorization URL.

```typescript
interface InitPaymentResponse {
  authorization_url: string;  // Redirect user here
  access_code: string;
  reference: string;          // Save this for verification
}

async function initializePayment(
  tierId: string,
  callbackUrl?: string
): Promise<InitPaymentResponse> {
  const token = localStorage.getItem('auth_token');

  const response = await fetch('/api/v2/payments/initialize', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({
      subscription_tier_id: tierId,
      callback_url: callbackUrl || `${window.location.origin}/subscription/callback`,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message);
  }

  return response.json().then(r => r.data);
}

// Usage in your subscribe handler
async function handleSubscribe(tierId: string) {
  try {
    // Show loading state
    setLoading(true);

    // Initialize payment
    const payment = await initializePayment(tierId);

    // Store reference for verification after redirect
    sessionStorage.setItem('payment_reference', payment.reference);

    // Redirect to Paystack payment page
    window.location.href = payment.authorization_url;

  } catch (error) {
    // Handle errors
    if (error.message.includes('already has an active subscription')) {
      toast.error('You already have an active subscription');
    } else {
      toast.error('Failed to initialize payment. Please try again.');
    }
  } finally {
    setLoading(false);
  }
}
```

---

#### Step 3: Handle Payment Callback

After payment, Paystack redirects the user back to your callback URL with query parameters.

```typescript
// pages/subscription/callback.tsx (Next.js example)

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

export default function PaymentCallback() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [status, setStatus] = useState<'verifying' | 'success' | 'failed'>('verifying');

  useEffect(() => {
    const reference = searchParams.get('reference') || searchParams.get('trxref');

    if (!reference) {
      setStatus('failed');
      return;
    }

    verifyPayment(reference);
  }, [searchParams]);

  async function verifyPayment(reference: string) {
    try {
      const token = localStorage.getItem('auth_token');

      const response = await fetch(
        `/api/v2/payments/verify?reference=${reference}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        }
      );

      const result = await response.json();

      if (response.ok && result.success) {
        setStatus('success');

        // Clear stored reference
        sessionStorage.removeItem('payment_reference');

        // Redirect to dashboard after short delay
        setTimeout(() => {
          router.push('/dashboard?subscription=activated');
        }, 2000);

      } else {
        setStatus('failed');
      }

    } catch (error) {
      console.error('Payment verification failed:', error);
      setStatus('failed');
    }
  }

  return (
    <div className="payment-callback">
      {status === 'verifying' && (
        <div>
          <Spinner />
          <p>Verifying your payment...</p>
        </div>
      )}

      {status === 'success' && (
        <div>
          <CheckIcon />
          <h2>Payment Successful!</h2>
          <p>Your subscription has been activated.</p>
          <p>Redirecting to dashboard...</p>
        </div>
      )}

      {status === 'failed' && (
        <div>
          <XIcon />
          <h2>Payment Failed</h2>
          <p>Something went wrong with your payment.</p>
          <button onClick={() => router.push('/pricing')}>
            Try Again
          </button>
        </div>
      )}
    </div>
  );
}
```

---

#### Step 4: Check User Subscription Status

After login or on protected pages, check if user has an active subscription.

```typescript
interface UserSubscription {
  id: string;
  subscription_id: string;
  subscription: SubscriptionTier;
  status: 'active' | 'cancelled' | 'expired' | 'past_due';
  is_active: boolean;
  current_period_start: string;
  current_period_end: string;
  cancel_at_period_end: boolean;
}

async function getUserSubscription(): Promise<UserSubscription | null> {
  const token = localStorage.getItem('auth_token');

  const response = await fetch('/api/v2/user/subscriptions?limit=1', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  const result = await response.json();

  // Return the first active subscription or null
  const subscriptions = result.data.data;
  return subscriptions.find((s: UserSubscription) => s.status === 'active') || null;
}

// Usage in a React context
function SubscriptionProvider({ children }) {
  const [subscription, setSubscription] = useState<UserSubscription | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getUserSubscription()
      .then(setSubscription)
      .finally(() => setLoading(false));
  }, []);

  // Check if user has access to a feature
  const hasFeature = (feature: string) => {
    if (!subscription) return false;
    return subscription.subscription.features?.[feature] === true;
  };

  // Check subscription tier
  const hasTier = (tier: string) => {
    if (!subscription) return false;
    const tierOrder = ['free', 'basic', 'premium', 'business'];
    const userTierIndex = tierOrder.indexOf(subscription.subscription.tier);
    const requiredTierIndex = tierOrder.indexOf(tier);
    return userTierIndex >= requiredTierIndex;
  };

  return (
    <SubscriptionContext.Provider value={{ subscription, loading, hasFeature, hasTier }}>
      {children}
    </SubscriptionContext.Provider>
  );
}
```

---

#### Step 5: Cancel Subscription

Allow users to cancel their subscription.

```typescript
async function cancelSubscription(): Promise<void> {
  const token = localStorage.getItem('auth_token');

  const response = await fetch('/api/v2/payments/cancel', {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message);
  }
}

// Usage in settings page
function SubscriptionSettings() {
  const { subscription, refetch } = useSubscription();
  const [cancelling, setCancelling] = useState(false);

  async function handleCancel() {
    const confirmed = await confirm(
      'Are you sure you want to cancel your subscription? ' +
      'You will retain access until the end of your current billing period.'
    );

    if (!confirmed) return;

    try {
      setCancelling(true);
      await cancelSubscription();
      toast.success('Subscription cancelled. Access ends on ' +
        formatDate(subscription.current_period_end));
      refetch();
    } catch (error) {
      toast.error(error.message);
    } finally {
      setCancelling(false);
    }
  }

  if (!subscription) return null;

  return (
    <div className="subscription-settings">
      <h3>Current Plan: {subscription.subscription.name}</h3>
      <p>Status: {subscription.status}</p>
      <p>Renews: {formatDate(subscription.current_period_end)}</p>

      {subscription.status === 'active' && !subscription.cancel_at_period_end && (
        <button
          onClick={handleCancel}
          disabled={cancelling}
          className="danger"
        >
          {cancelling ? 'Cancelling...' : 'Cancel Subscription'}
        </button>
      )}

      {subscription.cancel_at_period_end && (
        <p className="warning">
          Your subscription will end on {formatDate(subscription.current_period_end)}
        </p>
      )}
    </div>
  );
}
```

---

#### Step 6: Using Paystack Inline (Alternative)

Instead of redirecting, you can use Paystack's inline popup for a smoother UX.

```html
<!-- Include Paystack inline script -->
<script src="https://js.paystack.co/v1/inline.js"></script>
```

```typescript
// Initialize with Paystack inline popup
async function handleSubscribeInline(tierId: string) {
  try {
    // Get payment initialization data
    const payment = await initializePayment(tierId);

    // Use Paystack inline popup instead of redirect
    const handler = PaystackPop.setup({
      key: 'pk_test_xxxxx',  // Your Paystack public key
      email: user.email,
      amount: tierPrice * 100,  // Amount in kobo
      currency: 'NGN',
      ref: payment.reference,

      callback: function(response) {
        // Payment completed, verify on backend
        verifyPayment(response.reference);
      },

      onClose: function() {
        // User closed popup without completing payment
        toast.info('Payment cancelled');
      },
    });

    handler.openIframe();

  } catch (error) {
    toast.error(error.message);
  }
}
```

---

### Webhook Handling

Webhooks provide a reliable backup for payment confirmation. The API automatically handles:

| Event | Action |
|-------|--------|
| `charge.success` | Activates subscription |
| `subscription.create` | Links Paystack subscription ID |
| `subscription.disable` | Marks subscription as cancelled |
| `invoice.payment_failed` | Marks subscription as past_due |

Configure your webhook URL in Paystack dashboard:
```
https://your-api-domain.com/api/v2/payments/webhook
```

---

### Error Handling

#### Common Error Responses

| Status | Message | Cause | Frontend Action |
|--------|---------|-------|-----------------|
| 400 | "user already has an active subscription" | User trying to subscribe twice | Show current subscription, offer upgrade |
| 400 | "subscription tier not found" | Invalid tier ID | Refresh tier list |
| 400 | "Payment not successful" | Payment failed at Paystack | Show retry option |
| 401 | "User not authenticated" | Missing/invalid token | Redirect to login |
| 500 | "Failed to activate subscription" | Server error | Show generic error, suggest retry |

#### Error Handling Example

```typescript
async function handlePaymentError(error: Error) {
  const message = error.message.toLowerCase();

  if (message.includes('already has an active subscription')) {
    // Redirect to manage subscription
    router.push('/settings/subscription');
    toast.info('You already have an active subscription');

  } else if (message.includes('not authenticated')) {
    // Redirect to login
    router.push('/login?redirect=/pricing');

  } else if (message.includes('payment not successful')) {
    // Offer retry
    toast.error('Payment failed. Please try again.');

  } else {
    // Generic error
    toast.error('Something went wrong. Please try again later.');
    console.error('Payment error:', error);
  }
}
```

---

### API Reference Summary

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/v2/subscriptions` | GET | No | List all subscription tiers |
| `/api/v2/subscriptions/:id` | GET | No | Get single tier details |
| `/api/v2/payments/initialize` | POST | Yes | Start payment flow |
| `/api/v2/payments/verify` | GET | Yes | Verify payment & activate |
| `/api/v2/payments/cancel` | DELETE | Yes | Cancel subscription |
| `/api/v2/payments/webhook` | POST | No* | Paystack webhook |
| `/api/v2/user/subscriptions` | GET | Yes | Get user's subscriptions |

*Webhook uses signature verification instead of JWT auth.

---

## License

MIT
