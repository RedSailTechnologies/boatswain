import { Component, Inject, OnInit } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Clipboard } from '@angular/cdk/clipboard';

@Component({
  selector: 'app-message-dialog',
  templateUrl: './message-dialog.component.html',
  styleUrls: ['./message-dialog.component.sass']
})
export class MessageDialogComponent implements OnInit {
  public message: string;
  public reason: string;
  public showCopyButton: boolean;

  constructor(public dialog: MatDialogRef<MessageDialogComponent>, @Inject(MAT_DIALOG_DATA) data, private clipboard: Clipboard) {
    this.message = data['message'];
    this.reason = data['reason'];
    this.showCopyButton = data['showCopy'];
  }

  ngOnInit(): void {
  }

  copy(string: string): void {
    this.clipboard.copy(string)
  }
}
