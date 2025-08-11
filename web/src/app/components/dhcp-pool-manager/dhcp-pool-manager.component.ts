import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ApiService } from '../../api.service';
@Component({
  selector: 'app-dhcp-pool-manager',
  templateUrl: './dhcp-pool-manager.component.html',
  styleUrls: ['./dhcp-pool-manager.component.scss'],
  standalone: false
})
export class DhcpPoolManagerComponent implements OnInit {
  @Input() pools: any[] = [];
  @Output() poolsChange = new EventEmitter<any[]>();
  errors: string[] = [];
  form: FormGroup;
  showPoolModalMode = false;
  editIndex: number | null = null;

  constructor(
    private fb: FormBuilder,     
    private apiService: ApiService
  ) {
    this.form = this.fb.group({
      net_address: ['', Validators.required],
      netmask: ['', Validators.required],
      name: ['', Validators.required],
      start_address: ['', Validators.required],
      end_address: ['', Validators.required],
      gateway: ['', Validators.required],
    });
  }

  ngOnInit(): void {}

  showPoolModal(mode: 'add' | 'edit', id?: number) {

    if (mode === 'add') {
      this.openAddModal();
    } else if (mode === 'edit' && typeof id === 'number') {
      const index = this.pools.findIndex(pool => pool.id === id);
      if (index !== -1) {
        this.openEditModal(index);
      }
    }
    
  }

  openAddModal() {
    this.form.reset();
    this.editIndex = null;
    this.showPoolModalMode = true;
  }

  openEditModal(index: number) {
    this.editIndex = index;
    this.form.patchValue(this.pools[index]);
    this.showPoolModalMode = true;
  }

  submit() {
    if (this.form.invalid) return;

    const data = {
      ...this.form.value,
      only_serve_reimage: true,
      lease_time: 7000,
    };


    this.apiService.addPool(data).subscribe((resp: any) => {
      if (resp.error) {
        this.errors = resp.error;
      }
      if (resp) {
        this.pools.push(resp);
        this.form.reset();
      }
    });

    this.showPoolModalMode = false;
    this.editIndex = null;
  }

  delete(index: number) {
    this.pools.splice(index, 1);
    this.poolsChange.emit(this.pools);
  }
}
