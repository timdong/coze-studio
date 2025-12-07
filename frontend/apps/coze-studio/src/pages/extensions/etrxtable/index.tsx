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
 * Multi-dimensional Table System Page
 */

import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { TableList, TableEditor, TableViewer } from '@coze-extensions/etrxtable';
import type { TableBody } from '@coze-extensions/etrxtable';

export function EtrxTablePage() {
  const { space_id } = useParams<{ space_id: string }>();
  const workspaceId = space_id ? parseInt(space_id, 10) : 0;
  const [selectedTable, setSelectedTable] = useState<TableBody | null>(null);
  const [viewMode, setViewMode] = useState<'list' | 'edit' | 'view'>('list');

  if (!workspaceId) {
    return <div>无效的工作空间 ID</div>;
  }

  const handleTableSelect = (table: TableBody) => {
    setSelectedTable(table);
    setViewMode('view');
  };

  const handleCreateTable = (table: TableBody) => {
    setSelectedTable(table);
    setViewMode('view');
  };

  const handleEdit = () => {
    setViewMode('edit');
  };

  const handleBack = () => {
    setSelectedTable(null);
    setViewMode('list');
  };

  if (viewMode === 'list') {
    return (
      <div className="etrxtable-page">
        <TableList
          workspaceId={workspaceId}
          onTableSelect={handleTableSelect}
          onCreateTable={handleCreateTable}
        />
      </div>
    );
  }

  if (viewMode === 'edit' && selectedTable) {
    return (
      <div className="etrxtable-page">
        <button onClick={handleBack}>返回列表</button>
        <TableEditor
          tableId={selectedTable.id}
          onSave={handleBack}
          onCancel={handleBack}
        />
      </div>
    );
  }

  if (viewMode === 'view' && selectedTable) {
    return (
      <div className="etrxtable-page">
        <button onClick={handleBack}>返回列表</button>
        <button onClick={handleEdit}>编辑</button>
        <TableViewer tableId={selectedTable.id} onEdit={handleEdit} />
      </div>
    );
  }

  return null;
}

