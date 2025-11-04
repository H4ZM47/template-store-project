import {
  S3Client,
  PutObjectCommand,
  GetObjectCommand,
  DeleteObjectCommand,
  HeadObjectCommand,
} from '@aws-sdk/client-s3';
import { getSignedUrl } from '@aws-sdk/s3-request-presigner';
import { config } from '../config/config';
import { logger } from '../utils/logger';
import { v4 as uuidv4 } from 'uuid';

export interface UploadParams {
  file: Buffer;
  fileName: string;
  contentType: string;
  folder?: string;
}

export interface UploadResult {
  url: string;
  key: string;
  size: number;
}

export class StorageService {
  private client: S3Client;
  private bucket: string;

  constructor() {
    this.client = new S3Client({
      region: config.aws.region,
      credentials: {
        accessKeyId: config.aws.accessKeyId,
        secretAccessKey: config.aws.secretAccessKey,
      },
    });
    this.bucket = config.aws.s3Bucket;
  }

  async upload(params: UploadParams): Promise<UploadResult> {
    try {
      const fileExtension = params.fileName.split('.').pop();
      const key = params.folder
        ? `${params.folder}/${uuidv4()}.${fileExtension}`
        : `${uuidv4()}.${fileExtension}`;

      const command = new PutObjectCommand({
        Bucket: this.bucket,
        Key: key,
        Body: params.file,
        ContentType: params.contentType,
      });

      await this.client.send(command);

      const url = `https://${this.bucket}.s3.${config.aws.region}.amazonaws.com/${key}`;

      logger.info('File uploaded to S3', { key, size: params.file.length });

      return {
        url,
        key,
        size: params.file.length,
      };
    } catch (error) {
      logger.error('S3 upload error:', error);
      throw error;
    }
  }

  async getSignedDownloadUrl(key: string, expiresIn: number = 3600): Promise<string> {
    try {
      const command = new GetObjectCommand({
        Bucket: this.bucket,
        Key: key,
      });

      const url = await getSignedUrl(this.client, command, { expiresIn });
      logger.debug('Generated signed download URL', { key, expiresIn });
      return url;
    } catch (error) {
      logger.error('S3 get signed URL error:', error);
      throw error;
    }
  }

  async delete(key: string): Promise<void> {
    try {
      const command = new DeleteObjectCommand({
        Bucket: this.bucket,
        Key: key,
      });

      await this.client.send(command);
      logger.info('File deleted from S3', { key });
    } catch (error) {
      logger.error('S3 delete error:', error);
      throw error;
    }
  }

  async fileExists(key: string): Promise<boolean> {
    try {
      const command = new HeadObjectCommand({
        Bucket: this.bucket,
        Key: key,
      });

      await this.client.send(command);
      return true;
    } catch (error) {
      return false;
    }
  }

  extractKeyFromUrl(url: string): string | null {
    try {
      const bucketUrl = `https://${this.bucket}.s3.${config.aws.region}.amazonaws.com/`;
      if (url.startsWith(bucketUrl)) {
        return url.substring(bucketUrl.length);
      }
      return null;
    } catch (error) {
      return null;
    }
  }
}
