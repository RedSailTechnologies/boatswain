import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DeploymentDialogComponent } from './deployment-dialog.component';

describe('DeploymentDialogComponent', () => {
  let component: DeploymentDialogComponent;
  let fixture: ComponentFixture<DeploymentDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DeploymentDialogComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(DeploymentDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
