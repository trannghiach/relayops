"use client";

import { api } from "@/lib/api";
import type { JobItem } from "@/lib/types";
import { Button, Input, Space, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import Link from "next/link";
import { useEffect, useState } from "react";

function statusColor(status: string) {
    switch (status) {
        case "succeeded":
            return "green";
        case "pending":
            return "blue";
        case "processing":
            return "orange";
        case "dead_lettered":
            return "red";
        case "failed":
            return "volcano";
        default:
            return "default";
    }
}

export default function JobsPage() {
    const [jobs, setJobs] = useState<JobItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [query, setQuery] = useState("");

    async function fetchJobs(q = "") {
        setLoading(true);

        try {
            const res = await api.getJobs(q);
            setJobs(res.data);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchJobs();
    }, []);

    const columns: ColumnsType<JobItem> = [
        {
            title: "Job Type",
            dataIndex: "job_type",
        },
        {
            title: "Channel",
            dataIndex: "channel",
        },
        {
            title: "Status",
            dataIndex: "status",
            render: (status) => <Tag color={statusColor(status)}>{status}</Tag>,
        },
        {
            title: "Executions",
            render: (_, row) => `${row.attempts} total`,
        },
        {
            title: "Created At",
            dataIndex: "created_at",
            render: (v) => new Date(v).toLocaleString(),
        },
        {
            title: "Action",
            render: (_, row) => (
                <Link href={`/jobs/${row.id}`}>
                    <Button size="small">View</Button>
                </Link>
            ),
        },
    ];

    return (
        <>
            <Space direction="vertical" size="middle" style={{ width: "100%" }}>
                <Typography.Title level={2}>Jobs</Typography.Title>

                <Input.Search
                    placeholder="Search jobs by type, channel, or status"
                    allowClear
                    enterButton="Search"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onSearch={(value) => fetchJobs(value)}
                    style={{ maxWidth: 520 }}
                />

                <Table
                    rowKey="id"
                    loading={loading}
                    columns={columns}
                    dataSource={jobs}
                    pagination={{ pageSize: 10 }}
                />
            </Space>
        </>
    );
}