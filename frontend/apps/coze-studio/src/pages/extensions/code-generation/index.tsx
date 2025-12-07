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
 * Code Generation Page
 */

import React, { useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { CodeGenForm } from '@coze-extensions/code-generation';

export function CodeGenerationPage() {
  const { space_id } = useParams<{ space_id: string }>();
  const [searchParams] = useSearchParams();
  const tableIdParam = searchParams.get('table_id');
  const tableId = tableIdParam ? parseInt(tableIdParam, 10) : 0;

  if (!tableId) {
    return (
      <div className="code-generation-page">
        <div className="error">
          <p>请先选择一个表格</p>
          <p>使用方式: /space/:space_id/code-generation?table_id=1</p>
        </div>
      </div>
    );
  }

  return (
    <div className="code-generation-page">
      <CodeGenForm
        tableId={tableId}
        onCodeGenerated={(code) => {
          console.log('Generated code:', code);
        }}
      />
    </div>
  );
}

