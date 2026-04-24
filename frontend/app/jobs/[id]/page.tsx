"use client";

import { api } from "@/lib/api";
import type { DeliveryAttempt, JobDetail } from "@/lib/types";
import { Button, Card, Descriptions, Space, Table, Tag, Typography, message } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useParams } from "next/navigation";
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

export default function JobDetailPage() {
    const params = useParams<{ id: string }>();
    const jobID = params.id;

    const [job, setJob] = useState<JobDetail | null>(null);
    const [attempts, setAttempts] = useState<DeliveryAttempt[]>([]);
    const [loading, setLoading] = useState(true);

    async function load() {
        setLoading(true);
        try {
            const [jobRes, attemptsRes] = await Promise.all([
                api.getJob(jobID),
                api.getJobAttempts(jobID),
            ]);

            setJob(jobRes.data);
            setAttempts(attemptsRes.data);
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        load();
    }, [jobID]);

    async function handleReplay() {
        await api.replayJob(jobID);
        message.success("Job replayed");
        await load();
    }

    async function handleRetry() {
        await api.retryJob(jobID);
        message.success("Job retried");
        await load();
    }

    const attemptColumns: ColumnsType<DeliveryAttempt> = [
        {
            title: "Attempt",
            dataIndex: "attempt_no",
        },
        {
            title: "Provider",
            dataIndex: "provider",
        },
        {
            title: "Status",
            dataIndex: "status",
            render: (status) => <Tag color={statusColor(status)}>{status}</Tag>,
        },
        {
            title: "Error",
            dataIndex: "error_message",
            render: (v) => v || "-",
        },
        {
            title: "Started At",
            dataIndex: "started_at",
            render: (v) => new Date(v).toLocaleString(),
        },
    ];

    if (!job) {
        return <Typography.Text>Loading...</Typography.Text>;
    }

    return (
        <Space orientation="vertical" size="large" style={{ width: "100%" }}>
            <Space align="center">
                <Typography.Title level={2} style={{ margin: 0 }}>
                    Job Detail
                </Typography.Title>

                <Tag color={statusColor(job.status)}>{job.status}</Tag>
            </Space>

            <Card>
                <Descriptions bordered column={2}>
                    <Descriptions.Item label="Job ID">{job.id}</Descriptions.Item>
                    <Descriptions.Item label="Event ID">{job.event_id}</Descriptions.Item>
                    <Descriptions.Item label="Job Type">{job.job_type}</Descriptions.Item>
                    <Descriptions.Item label="Channel">{job.channel}</Descriptions.Item>
                    <Descriptions.Item label="Attempts">
                        {job.attempts}/{job.max_attempts}
                    </Descriptions.Item>
                    <Descriptions.Item label="Last Error">
                        {job.last_error || "-"}
                    </Descriptions.Item>
                    <Descriptions.Item label="Created At">
                        {new Date(job.created_at).toLocaleString()}
                    </Descriptions.Item>
                    <Descriptions.Item label="Updated At">
                        {new Date(job.updated_at).toLocaleString()}
                    </Descriptions.Item>
                </Descriptions>
            </Card>

            <Space>
                <Button
                    onClick={handleRetry}
                    disabled={!["failed", "pending", "processing"].includes(job.status)}
                >
                    Retry
                </Button>

                <Button
                    danger
                    type="primary"
                    onClick={handleReplay}
                    disabled={job.status !== "dead_lettered"}
                >
                    Replay Dead Letter
                </Button>
            </Space>

            <Card title="Payload">
                <pre>{JSON.stringify(job.payload, null, 2)}</pre>
            </Card>

            <Card title="Delivery Attempts">
                <Table
                    rowKey="id"
                    loading={loading}
                    columns={attemptColumns}
                    dataSource={attempts}
                    pagination={false}
                />
            </Card>
        </Space>
    );
}