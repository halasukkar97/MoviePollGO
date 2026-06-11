import { useEffect, useMemo, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { apiClient } from '../../api/client';
import type { TranslationKey } from '../../i18n/useTranslations';
import './Breadcrumbs.scss';

interface BreadcrumbsProps {
  t: (key: TranslationKey) => string;
}

type BreadcrumbItem = {
  label: string;
  to?: string;
};

function getPollCodeFromPath(pathname: string) {
  const match = pathname.match(/^\/polls\/([^/]+)/);
  return match?.[1] === 'new' ? '' : match?.[1] ?? '';
}

function getPollCodeFromSearch(search: string) {
  return new URLSearchParams(search).get('pollCode') ?? '';
}

// Breadcrumbs builds page location links from the current React Router path.
export function Breadcrumbs({ t }: BreadcrumbsProps) {
  const location = useLocation();
  const [pollName, setPollName] = useState('');

  const pathname = location.pathname;
  const pollCode = getPollCodeFromPath(pathname) || getPollCodeFromSearch(location.search);

  // Poll routes use the public pollCode in the URL, then load the readable poll name.
  useEffect(() => {
    let isMounted = true;

    async function loadPollName() {
      if (!pollCode) {
        setPollName('');
        return;
      }

      try {
        const poll = await apiClient.getPoll(pollCode);
        if (isMounted) {
          setPollName(poll.name || pollCode);
        }
      } catch {
        if (isMounted) {
          setPollName(pollCode);
        }
      }
    }

    loadPollName();

    return () => {
      isMounted = false;
    };
  }, [pollCode]);

  const items = useMemo<BreadcrumbItem[]>(() => {
    const home = { label: t('nav.home'), to: '/' };
    const history = { label: t('nav.history'), to: '/history' };

    if (pathname === '/') {
      return [{ label: t('nav.home') }];
    }

    if (pathname === '/history') {
      return [home, { label: t('nav.history') }];
    }

    if (pathname === '/polls/new') {
      return [home, { label: t('nav.createPoll') }];
    }

    if (pathname === '/join') {
      return [home, { label: t('nav.joinPoll') }];
    }

    if (pathname.startsWith('/polls/')) {
      return [home, history, { label: pollName || pollCode }];
    }

    if (pathname === '/results') {
      return pollCode
        ? [home, history, { label: pollName || pollCode, to: '/polls/' + pollCode }, { label: t('results.title') }]
        : [home, history, { label: t('results.title') }];
    }

    return [home, { label: pathname.replace('/', '') || t('nav.home') }];
  }, [pathname, pollCode, pollName, t]);

  return (
    <nav className="breadcrumbs" aria-label="Breadcrumb">
      {items.map((item, index) => {
        const isLast = index === items.length - 1;

        return (
          <span className="breadcrumb-item" key={item.label + index}>
            {item.to && !isLast ? <Link to={item.to}>{item.label}</Link> : <span>{item.label}</span>}
            {!isLast ? <span className="breadcrumb-separator">/</span> : null}
          </span>
        );
      })}
    </nav>
  );
}
