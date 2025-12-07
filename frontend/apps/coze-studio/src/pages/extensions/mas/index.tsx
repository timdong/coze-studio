/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * Multi-Agent System Page
 */

import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { AgentList, SessionList, TaskList } from '@coze-extensions/mas';
import type { Agent, MASSession } from '@coze-extensions/mas';

export function MASPage() {
  const { space_id } = useParams<{ space_id: string }>();
  const workspaceId = space_id ? parseInt(space_id, 10) : 0;
  const [activeTab, setActiveTab] = useState<'agents' | 'sessions' | 'tasks'>('agents');
  const [selectedSession, setSelectedSession] = useState<MASSession | null>(null);

  if (!workspaceId) {
    return <div>无效的工作空间 ID</div>;
  }

  const handleSessionSelect = (session: MASSession) => {
    setSelectedSession(session);
    setActiveTab('tasks');
  };

  if (activeTab === 'tasks' && selectedSession) {
    return (
      <div className="mas-page">
        <div className="mas-tabs">
          <button onClick={() => setActiveTab('agents')}>Agents</button>
          <button onClick={() => setActiveTab('sessions')}>Sessions</button>
          <button onClick={() => setActiveTab('tasks')} className="active">
            Tasks
          </button>
        </div>
        <button onClick={() => setSelectedSession(null)}>返回 Sessions</button>
        <TaskList sessionId={selectedSession.id} />
      </div>
    );
  }

  return (
    <div className="mas-page">
      <div className="mas-tabs">
        <button
          onClick={() => {
            setActiveTab('agents');
            setSelectedSession(null);
          }}
          className={activeTab === 'agents' ? 'active' : ''}
        >
          Agents
        </button>
        <button
          onClick={() => {
            setActiveTab('sessions');
            setSelectedSession(null);
          }}
          className={activeTab === 'sessions' ? 'active' : ''}
        >
          Sessions
        </button>
      </div>

      {activeTab === 'agents' && (
        <AgentList workspaceId={workspaceId} />
      )}

      {activeTab === 'sessions' && (
        <SessionList
          workspaceId={workspaceId}
          onSessionSelect={handleSessionSelect}
        />
      )}
    </div>
  );
}

